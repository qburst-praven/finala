package meilisearch

import (
	"encoding/json"
	"errors"
	"finala/api/config"
	"finala/api/storage"
	"finala/interpolation"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidQuery            = errors.New("invalid query")
	ErrAggregationTermNotFound = errors.New("aggregation terms was not found")
)

const (
	// prefixDayIndex defines the index name of the current day
	prefixIndexName = "finala-%s"
)

// StorageManager describes meilisearchStorage
type StorageManager struct {
	client          meilisearchDescriptor
	currentIndexDay string
}

// NewStorageManager creates new Meilisearch storage
func NewStorageManager(conf config.MeilisearchConfig) (*StorageManager, error) {
	client, err := NewClient(conf)
	if err != nil {
		return nil, err
	}

	storageManager := &StorageManager{
		client: client,
	}

	if !storageManager.setCreateCurrentIndexDay() {
		return nil, errors.New("could not create index")
	}

	go func() {
		for {
			now := time.Now().In(time.UTC)
			diff := storageManager.getDurationUntilTomorrow(now)
			log.WithFields(log.Fields{
				"now":      now,
				"duration": diff,
			}).Info("change index in")
			<-time.After(diff)
			storageManager.setCreateCurrentIndexDay()
		}
	}()

	return storageManager, nil
}

// Save new documents
func (sm *StorageManager) Save(data string) bool {
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(data), &doc); err != nil {
		log.WithError(err).Error("Failed to unmarshal document")
		return false
	}

	// Add an ID field if not present (required by Meilisearch)
	if _, ok := doc["id"]; !ok {
		doc["id"] = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	err := sm.client.Index(sm.currentIndexDay, doc)
	if err != nil {
		log.WithFields(log.Fields{
			"index": sm.currentIndexDay,
			"data":  data,
		}).WithError(err).Error("Fail to save document")
		return false
	}

	return true
}

// GetSummary returns executions summary
func (sm *StorageManager) GetSummary(executionID string, filters map[string]string) (map[string]storage.CollectorsSummary, error) {
	summary := make(map[string]storage.CollectorsSummary)

	// 1. Fetch and process service_status events for status and error messages
	serviceStatusEvents, err := sm.client.Search(sm.currentIndexDay, map[string]interface{}{
		"q":         "",
		"filter_by": fmt.Sprintf("EventType=service_status AND ExecutionID=%s", executionID),
		// Potentially add limit if there can be many status events per service, though unlikely for summary.
		// Default MeiliSearch limit is 20, might need to be higher if many resource types.
		// For now, assuming default limit is sufficient or all relevant statuses are captured.
	})
	if err != nil {
		log.WithError(err).Error("error when trying to get service_status summary data")
		// Decide if we should return partial data or error out. For now, continue to try fetching resource data.
	}

	if serviceStatusEvents != nil {
		for _, hit := range serviceStatusEvents.Hits {
			var statusData storage.Summary // storage.Summary is from api/storage/structs.go
			hitData, err := json.Marshal(hit)
			if err != nil {
				log.WithError(err).Error("could not marshal service_status hit")
				continue
			}
			if err := json.Unmarshal(hitData, &statusData); err != nil {
				log.WithError(err).Error("could not parse service_status row")
				continue
			}

			// Ensure we only keep the latest status event if multiple exist (though GetSummary implies one summary)
			existing, found := summary[statusData.ResourceName]
			if found && statusData.EventTime < existing.EventTime {
				continue
			}

			summary[statusData.ResourceName] = storage.CollectorsSummary{
				ResourceName: statusData.ResourceName,
				Status:       statusData.Data.Status,
				ErrorMessage: statusData.Data.ErrorMessage,
				EventTime:    statusData.EventTime,
				// ResourceCount and TotalSpent will be populated from resource_detected events
			}
		}
	}

	// 2. Fetch and process resource_detected events for costs and counts
	resourceDetectedEvents, err := sm.client.Search(sm.currentIndexDay, map[string]interface{}{
		"q":         "",
		"filter_by": fmt.Sprintf("EventType=resource_detected AND ExecutionID=%s", executionID),
		"limit":     1000, // Assuming up to 1000 resources per execution for summary. Adjust if necessary.
	})

	if err != nil {
		log.WithError(err).Error("error when trying to get resource_detected summary data")
		// If we can't get resource data, the summary will be incomplete (costs will be 0).
		// Return the summary populated so far (with statuses) or error out.
		// For now, return what we have, which might be just statuses.
		return summary, err // Or, if statuses are primary, return summary, nil
	}

	if resourceDetectedEvents != nil {
		for _, hit := range resourceDetectedEvents.Hits {
			var eventDataMap map[string]interface{}
			hitData, err := json.Marshal(hit)
			if err != nil {
				log.WithError(err).Error("could not marshal resource_detected hit")
				continue
			}
			if err := json.Unmarshal(hitData, &eventDataMap); err != nil {
				log.WithError(err).Error("could not unmarshal resource_detected hit to map")
				continue
			}

			resourceName, rnOK := eventDataMap["ResourceName"].(string)
			if !rnOK {
				log.Error("ResourceName missing or not a string in resource_detected event")
				continue
			}

			dataField, dataOK := eventDataMap["Data"].(map[string]interface{})
			if !dataOK {
				log.WithField("resourceName", resourceName).Error("Data field missing or not a map in resource_detected event")
				continue
			}

			// PriceDetectedFields are embedded, so PricePerMonth should be a top-level field in Data
			pricePerMonth, ppmOK := dataField["PricePerMonth"].(float64)
			if !ppmOK {
				// Some resources like Lambda might not have PricePerMonth directly.
				// Handle this gracefully, e.g. by attempting to get PricePerHour or setting to 0.
				// For now, we log and skip if PricePerMonth is not found or not a float64.
				// This might be an area for improvement based on how various resources report costs.
				log.WithFields(log.Fields{
					"resourceName": resourceName,
					"dataField":    dataField,
				}).Warn("PricePerMonth not found or not a float64 in resource_detected event Data")
				// Continue accumulating other data, or decide if this resource should have a $0 cost.
			}

			currentSummary := summary[resourceName] // Get existing summary (could be just status info)
			currentSummary.ResourceName = resourceName // Ensure ResourceName is set
			currentSummary.ResourceCount++
			currentSummary.TotalSpent += pricePerMonth
			// Status, ErrorMessage, EventTime are already set from service_status or will be default if no status event.

			summary[resourceName] = currentSummary
		}
	}

	// Fill in any missing ResourceNames for services that had detected resources but no explicit status event
	// (though typically a CollectStart/Finish should exist)
	if resourceDetectedEvents != nil {
		for _, hit := range resourceDetectedEvents.Hits {
				var eventDataMap map[string]interface{}
				hitData, err := json.Marshal(hit)
				if err != nil {continue}
				if err := json.Unmarshal(hitData, &eventDataMap); err != nil {continue}
				resourceName, rnOK := eventDataMap["ResourceName"].(string)
				if !rnOK {continue}

				if _, exists := summary[resourceName]; !exists {
						summary[resourceName] = storage.CollectorsSummary{
								ResourceName: resourceName,
								// Status, ErrorMessage, EventTime will be zero/empty if no corresponding service_status event
						}
				}
		}
	}

	return summary, nil
}

// GetExecutions returns list of executions
func (sm *StorageManager) GetExecutions(queryLimit int) ([]storage.Executions, error) {
	executions := []storage.Executions{}

	searchParams := map[string]interface{}{
		"q": "",
		"filter_by": "EventType=service_status",
	}

	result, err := sm.client.Search(sm.currentIndexDay, searchParams)
	if err != nil {
		log.WithError(err).Error("error when trying to get executions collectors")
		return executions, ErrInvalidQuery
	}

	// Group by ExecutionID manually since Meilisearch doesn't support group by
	executionMap := make(map[string]bool)
	for _, hit := range result.Hits {
		var execData struct {
			ExecutionID string `json:"ExecutionID"`
		}
		hitData, err := json.Marshal(hit)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(hitData, &execData); err != nil {
			continue
		}

		if _, exists := executionMap[execData.ExecutionID]; !exists {
			timestamp, err := interpolation.ExtractTimestamp(execData.ExecutionID)
			if err != nil {
				timestamp = 0
			}

			// Use the correct fields: ID instead of ExecutionID, and set Time
			executions = append(executions, storage.Executions{
				ID:   execData.ExecutionID,
				Name: "Execution " + execData.ExecutionID,
				Time: time.Unix(timestamp, 0),
			})
			executionMap[execData.ExecutionID] = true
		}
	}

	return executions, nil
}

// GetResources returns list of resources
func (sm *StorageManager) GetResources(resourceType string, executionID string, filters map[string]string, search string) ([]map[string]interface{}, error) {
	var resources []map[string]interface{}

	searchQuery := ""
	if search != "" {
		searchQuery = search
	}

	searchParams := map[string]interface{}{
		"q": searchQuery,
		"filter_by": fmt.Sprintf("EventType=resource_detected AND ExecutionID=%s AND ResourceName=%s", executionID, resourceType),
	}

	result, err := sm.client.Search(sm.currentIndexDay, searchParams)
	if err != nil {
		log.WithError(err).Error("meilisearch query error")
		return resources, err
	}

	for _, hit := range result.Hits {
		rowData := make(map[string]interface{})
		hitData, err := json.Marshal(hit)
		if err != nil {
			log.WithError(err).Error("error when trying to marshal document")
			continue
		}
		if err := json.Unmarshal(hitData, &rowData); err != nil {
			log.WithError(err).Error("error when trying to parse search result hits data")
			continue
		}

		resources = append(resources, rowData)
	}

	return resources, nil
}

// GetResourceTrends returns resource trends
func (sm *StorageManager) GetResourceTrends(resourceType string, filters map[string]string, limit int) ([]storage.ExecutionCost, error) {
	var resources []storage.ExecutionCost

	// Build filter string for Meilisearch
	filterStr := fmt.Sprintf("ResourceName=%s AND EventType!=service_status", resourceType)
	
	// Add additional filters if any
	for key, value := range filters {
		filterStr += fmt.Sprintf(" AND %s=%s", key, value)
	}

	searchParams := map[string]interface{}{
		"q": "",
		"filter_by": filterStr,
	}

	result, err := sm.client.Search(sm.currentIndexDay, searchParams)
	if err != nil {
		log.WithError(err).Error("meilisearch query error")
		return resources, err
	}

	// Group by ExecutionID manually since Meilisearch doesn't support group by
	executionCosts := make(map[string]float64)
	for _, hit := range result.Hits {
		var execData struct {
			ExecutionID string                 `json:"ExecutionID"`
			Data        map[string]interface{} `json:"Data"`
		}
		hitData, err := json.Marshal(hit)
		if err != nil {
			log.WithError(err).Error("Error marshaling document hit")
			continue
		}
		if err := json.Unmarshal(hitData, &execData); err != nil {
			log.WithError(err).Error("Error unmarshaling data")
			continue
		}

		if priceData, ok := execData.Data["PricePerMonth"]; ok {
			if price, ok := priceData.(float64); ok {
				executionCosts[execData.ExecutionID] += price
			} else {
				// Handle case where PricePerMonth might be a string or other type
				priceStr, ok := priceData.(string)
				if ok {
					if priceVal, err := fmt.Sscanf(priceStr, "%f", new(float64)); err == nil {
						executionCosts[execData.ExecutionID] += float64(priceVal)
					}
				}
			}
		}
	}

	// Convert to array and sort by timestamp
	for execID, costSum := range executionCosts {
		timestamp, err := interpolation.ExtractTimestamp(execID)
		if err != nil {
			timestamp = 0
		}

		resources = append(resources, storage.ExecutionCost{
			ExecutionID:        execID,
			ExtractedTimestamp: timestamp,
			CostSum:            costSum,
		})
	}

	return resources, nil
}

// GetExecutionTags returns execution tags
func (sm *StorageManager) GetExecutionTags(executionID string) (map[string][]string, error) {
	tags := map[string][]string{}

	searchParams := map[string]interface{}{
		"q": "",
		"filter_by": fmt.Sprintf("EventType=resource_detected AND ExecutionID=%s", executionID),
	}

	result, err := sm.client.Search(sm.currentIndexDay, searchParams)
	if err != nil {
		log.WithError(err).Error("got a meilisearch error while running the query")
		return tags, err
	}

	log.WithFields(log.Fields{
		"hits_count": len(result.Hits),
		"execution_id": executionID,
	}).Debug("Processing tags from search results")

	for _, hit := range result.Hits {
		// First try to unmarshal with expected structure
		var tagsData struct {
			Data struct {
				Tag map[string]string `json:"Tag"`
			} `json:"Data"`
		}

		hitData, err := json.Marshal(hit)
		if err != nil {
			log.WithError(err).Debug("Error marshaling document hit")
			continue
		}
		
		if err := json.Unmarshal(hitData, &tagsData); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"hit": string(hitData),
			}).Debug("Error parsing tags structure, trying alternative structure")
			
			// Try alternative structure where Data.Tag might be a generic map
			var altTagsData struct {
				Data map[string]interface{} `json:"Data"`
			}
			
			if err := json.Unmarshal(hitData, &altTagsData); err != nil {
				log.WithError(err).Debug("Error parsing alternative tags structure")
				continue
			}
			
			// Check if Tag exists in Data map
			if tagData, ok := altTagsData.Data["Tag"]; ok {
				// Try to cast to map[string]string or map[string]interface{}
				if tagMap, ok := tagData.(map[string]string); ok {
					for key, value := range tagMap {
						tags[key] = append(tags[key], value)
					}
				} else if tagMapIface, ok := tagData.(map[string]interface{}); ok {
					for key, value := range tagMapIface {
						if strValue, ok := value.(string); ok {
							tags[key] = append(tags[key], strValue)
						}
					}
				}
			}
			
			continue
		}

		// If we got here, we successfully parsed the expected structure
		for key, value := range tagsData.Data.Tag {
			tags[key] = append(tags[key], value)
		}
	}

	// Make sure the values of each tag unique
	for tagName, tagValues := range tags {
		tags[tagName] = interpolation.UniqueStr(tagValues)
	}

	return tags, nil
}

// createIndex creates a new index if it doesn't exist
func (sm *StorageManager) createIndex(name string) error {
	exists, err := sm.client.IndexExists(name)
	if err != nil {
		log.WithFields(log.Fields{
			"index": name,
		}).WithError(err).Error("Error when trying to check if meilisearch index exists")
		return err
	}

	if exists {
		log.WithField("index", name).Info("index already exists")
		return nil
	}

	err = sm.client.CreateIndex(name)
	if err != nil {
		log.WithFields(log.Fields{
			"index": name,
		}).WithError(err).Error("Error when trying to create meilisearch index")
		return err
	}

	log.WithField("index", name).Info("index created successfully")
	return nil
}

// getDurationUntilTomorrow returns the duration time until tomorrow
func (sm *StorageManager) getDurationUntilTomorrow(now time.Time) time.Duration {
	zone, _ := now.Zone()
	location, err := time.LoadLocation(zone)
	if err != nil {
		log.WithError(err).WithField("zone", zone).Warn("zone name not found")
		location = time.UTC
	}

	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
	return tomorrow.Sub(now)
}

// setCreateCurrentIndexDay create and set the current day as index
func (sm *StorageManager) setCreateCurrentIndexDay() bool {
	dt := time.Now().In(time.UTC)
	newIndex := fmt.Sprintf(prefixIndexName, dt.Format("01-02-2006"))
	log.WithFields(log.Fields{
		"current_index_day":    sm.currentIndexDay,
		"to_current_index_day": newIndex,
	}).Info("change current index day")
	err := sm.createIndex(newIndex)
	if err != nil {
		return false
	}

	sm.currentIndexDay = newIndex
	return true
} 