package meilisearch

import (
	"bytes"
	"encoding/json"
	"finala/api/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

const (
	// connectionInterval defines the time duration to wait until the next connection retry
	connectionInterval = 5 * time.Second
	// connectionTimeout defines the maximum time duration until the API returns a connection error
	connectionTimeout = 60 * time.Second
)

// meilisearchDescriptor is the Meilisearch root interface that matches ES functionality
type meilisearchDescriptor interface {
	Index(index string, document interface{}) error
	Search(index string, query interface{}) (*meilisearch.SearchResponse, error)
	CreateIndex(name string) error
	IndexExists(name string) (bool, error)
}

// meilisearchClient implements the meilisearchDescriptor interface
type meilisearchClient struct {
	client *meilisearch.Client
	host   string
	apiKey string
}

// NewClient creates new Meilisearch client
func NewClient(conf config.MeilisearchConfig) (*meilisearchClient, error) {
	var client *meilisearchClient
	c := make(chan int, 1)
	var err error

	go func() {
		for {
			client, err = getMeilisearchClient(conf)
			if err == nil {
				// Test connection with a simple health check
				err = client.healthCheck()
				if err == nil {
					break
				}
			}
			log.WithFields(log.Fields{
				"endpoint": conf.Endpoints,
			}).WithError(err).Warn(fmt.Sprintf("could not initialize connection to Meilisearch, retrying in %v", connectionInterval))
			time.Sleep(connectionInterval)
		}
		c <- 1
	}()

	select {
	case <-c:
	case <-time.After(connectionTimeout):
		err = fmt.Errorf("could not connect Meilisearch, timed out after %v", connectionTimeout)
		log.WithError(err).Error("connection Error")
	}

	return client, err
}

// getMeilisearchClient creates new Meilisearch client
func getMeilisearchClient(conf config.MeilisearchConfig) (*meilisearchClient, error) {
	log.Infof("Creating Meilisearch client with endpoints: %v", conf.Endpoints)
	
	// For v0.12.0, set the apiKey directly and not through the config
	// since many methods might not be using the config correctly
	host := conf.Endpoints[0]
	apiKey := conf.Password
	
	// Create client with empty config, we'll manually add Authorization headers
	config := meilisearch.Config{
		Host:   host,
		APIKey: apiKey,
	}
	client := meilisearch.NewClient(config)

	return &meilisearchClient{
		client: client,
		host:   host,
		apiKey: apiKey,
	}, nil
}

// healthCheck verifies the Meilisearch connection
func (m *meilisearchClient) healthCheck() error {
	// For v0.12.0, we'll make a direct HTTP request to check health
	req, err := http.NewRequest("GET", m.host+"/health", nil)
	if err != nil {
		return err
	}
	
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("health check failed with status: %s, body: %s", resp.Status, string(body))
	}
	
	log.Info("Successfully connected to Meilisearch")
	return nil
}

// Index implements document indexing
func (m *meilisearchClient) Index(index string, document interface{}) error {
	url := fmt.Sprintf("%s/indexes/%s/documents", m.host, index)
	
	jsonDoc, err := json.Marshal(document)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonDoc))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Failed to index document: status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// Search implements search functionality
func (m *meilisearchClient) Search(index string, query interface{}) (*meilisearch.SearchResponse, error) {
	searchParams := query.(map[string]interface{})
	
	// Build the search query parameters
	q := ""
	if query, ok := searchParams["q"].(string); ok {
		q = query
	}
	
	// Add filter if present
	filter := ""
	if filterVal, ok := searchParams["filter_by"].(string); ok && filterVal != "" {
		filter = strings.ReplaceAll(filterVal, " = ", "=")
		filter = strings.ReplaceAll(filter, " != ", "!=")
	}
	
	// Build the search request URL with query parameters
	url := fmt.Sprintf("%s/indexes/%s/search", m.host, index)
	
	// Create the request JSON payload
	requestBody := map[string]interface{}{
		"q": q,
		"limit": 1000,
	}
	
	if filter != "" {
		requestBody["filter"] = filter
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	
	log.Debugf("Searching index %s with query: %s, filter: %s", index, q, filter)
	
	// Implementation with retry logic for filterability errors
	maxRetries := 3
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Create the HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		
		req.Header.Set("Content-Type", "application/json")
		if m.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+m.apiKey)
		}
		
		// Execute the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"index": index,
				"query": q,
				"filter": filter,
				"error": err,
				"attempt": attempt,
			}).Error("Error sending search request")
			lastErr = err
			continue
		}
		
		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			log.WithError(err).Error("Error reading response body")
			lastErr = err
			continue
		}
		
		// Handle error responses
		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("Search failed: status: %d, body: %s", resp.StatusCode, string(body))
			log.WithFields(log.Fields{
				"status": resp.StatusCode,
				"body": string(body),
				"attempt": attempt,
			}).Error("Search request failed")
			
			// If this is a filterable attributes error, try to fix it
			if resp.StatusCode == http.StatusBadRequest && 
			   strings.Contains(string(body), "not filterable") && 
			   strings.Contains(string(body), "does not have configured filterable attributes") {
				
				log.Warnf("Filterable attributes not properly configured for index %s, attempting to fix...", index)
				
				// Try to re-configure the index settings
				fixErr := m.configureIndexSettings(index)
				if fixErr != nil {
					log.WithError(fixErr).Warn("Failed to fix index settings")
				} else {
					log.Info("Successfully reconfigured index settings, retrying search")
				}
				
				// Wait before retrying
				waitTime := time.Duration(attempt*2) * time.Second
				time.Sleep(waitTime)
				lastErr = fmt.Errorf(errMsg)
				continue
			}
			
			return nil, fmt.Errorf(errMsg)
		}
		
		// Parse the response
		var searchResult meilisearch.SearchResponse
		err = json.Unmarshal(body, &searchResult)
		if err != nil {
			lastErr = err
			continue
		}
		
		return &searchResult, nil
	}
	
	return nil, fmt.Errorf("search failed after %d attempts: %v", maxRetries, lastErr)
}

// CreateIndex creates a new index
func (m *meilisearchClient) CreateIndex(name string) error {
	// Check if index exists first
	exists, err := m.IndexExists(name)
	if err != nil {
		log.WithFields(log.Fields{
			"index": name,
			"error": err,
		}).Error("Error checking if index exists")
		return err
	}
	
	if exists {
		log.Infof("Index %s already exists", name)
		return m.configureIndexSettings(name)
	}
	
	// Create a new index using direct HTTP request
	url := fmt.Sprintf("%s/indexes", m.host)
	payload := map[string]string{
		"uid":        name,
		"primaryKey": "id",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	
	log.Infof("Creating index: %s with URL: %s", name, url)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"index": name,
			"error": err,
		}).Error("Error sending request to create index")
		return err
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Failed to create index: %s, status: %d, body: %s", name, resp.StatusCode, string(body))
	}
	
	log.Infof("Index %s created successfully with response: %s", name, string(body))
	
	// Wait for index creation to complete
	time.Sleep(2 * time.Second)
	
	return m.configureIndexSettings(name)
}

// configureIndexSettings sets up the necessary index settings for optimal search
func (m *meilisearchClient) configureIndexSettings(name string) error {
	log.Infof("Configuring settings for index %s", name)
	
	// Define all settings
	settings := map[string]interface{}{
		"searchableAttributes": []string{
			"ResourceName",
			"ExecutionID",
			"EventType",
			"Data",
			"*",  // Make all fields searchable
		},
		"rankingRules": []string{
			"words",
			"typo",
			"proximity",
			"attribute",
			"sort",
			"Timestamp:desc",
			"exactness",
		},
		"filterableAttributes": []string{
			"ResourceName",
			"ExecutionID",
			"EventType",
			"Timestamp",
			"id",
		},
		"sortableAttributes": []string{
			"Timestamp",
		},
		"faceting": map[string]interface{}{
			"maxValuesPerFacet": 100,
		},
	}
	
	jsonSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	
	// Try up to 5 times with increasing delays
	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Infof("Attempting to configure index settings (attempt %d/%d)", attempt, maxRetries)
		
		// Apply the settings
		settingsUrl := fmt.Sprintf("%s/indexes/%s/settings", m.host, name)
		req, err := http.NewRequest("PATCH", settingsUrl, bytes.NewBuffer(jsonSettings))
		if err != nil {
			return err
		}
		
		req.Header.Set("Content-Type", "application/json")
		if m.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+m.apiKey)
		}
		
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"index": name,
				"error": err,
				"attempt": attempt,
			}).Warn("Error sending request to configure index settings")
			
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Infof("Retrying in %v", waitTime)
				time.Sleep(waitTime)
				continue
			}
			return err
		}
		
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			log.WithFields(log.Fields{
				"index": name,
				"status": resp.StatusCode,
				"response": string(body),
				"attempt": attempt,
			}).Warn("Failed to configure index settings")
			
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Infof("Retrying in %v", waitTime)
				time.Sleep(waitTime)
				continue
			}
			return fmt.Errorf("Failed to configure settings: status: %d, body: %s", resp.StatusCode, string(body))
		}
		
		log.Infof("Settings update for index %s accepted, waiting for it to be processed...", name)
		
		// Wait for settings to be applied - increase with each attempt
		waitTime := time.Duration(attempt*3) * time.Second
		time.Sleep(waitTime)
		
		// Verify settings were actually applied by checking the filterableAttributes
		// This is the most critical setting for our application
		verified := false
		for verifyAttempt := 1; verifyAttempt <= 3; verifyAttempt++ {
			filterableUrl := fmt.Sprintf("%s/indexes/%s/settings/filterable-attributes", m.host, name)
			verifyReq, err := http.NewRequest("GET", filterableUrl, nil)
			if err != nil {
				log.Warn("Could not create request to verify settings")
				break
			}
			
			if m.apiKey != "" {
				verifyReq.Header.Set("Authorization", "Bearer "+m.apiKey)
			}
			
			verifyResp, err := http.DefaultClient.Do(verifyReq)
			if err != nil {
				log.Warn("Could not verify filterable attributes, will retry")
				time.Sleep(time.Duration(verifyAttempt) * time.Second)
				continue
			}
			
			verifyBody, _ := ioutil.ReadAll(verifyResp.Body)
			verifyResp.Body.Close()
			
			if verifyResp.StatusCode == http.StatusOK {
				// Check if EventType is in the filterable attributes
				var filterableAttrs []string
				err = json.Unmarshal(verifyBody, &filterableAttrs)
				if err != nil {
					log.Warnf("Could not parse filterable attributes response: %s", string(verifyBody))
					time.Sleep(time.Duration(verifyAttempt) * time.Second)
					continue
				}
				
				// Check if critical attributes are in the list
				eventTypeFound := false
				for _, attr := range filterableAttrs {
					if attr == "EventType" {
						eventTypeFound = true
						break
					}
				}
				
				if eventTypeFound {
					log.Infof("Successfully verified filterable attributes for index %s: %s", name, string(verifyBody))
					verified = true
					break
				} else {
					log.Warnf("EventType not found in filterable attributes: %s", string(verifyBody))
				}
			} else {
				log.Warnf("Failed to verify filterable attributes, status: %d, body: %s", 
					verifyResp.StatusCode, string(verifyBody))
			}
			
			time.Sleep(time.Duration(verifyAttempt) * time.Second)
		}
		
		if verified {
			log.Infof("Successfully configured and verified index settings for %s", name)
			return nil
		}
		
		log.Warnf("Could not verify settings were properly applied for index %s, retrying full configuration", name)
	}
	
	return fmt.Errorf("Failed to configure and verify settings after %d attempts", maxRetries)
}

// IndexExists checks if an index exists
func (m *meilisearchClient) IndexExists(name string) (bool, error) {
	// For v1.2, make a direct HTTP request to check if index exists
	url := fmt.Sprintf("%s/indexes/%s", m.host, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	// If the status code is 200, the index exists
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	
	// If the status code is 404, the index doesn't exist
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	
	// Any other status code is an error
	body, _ := ioutil.ReadAll(resp.Body)
	return false, fmt.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
} 