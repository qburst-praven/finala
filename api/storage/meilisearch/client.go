package meilisearch

import (
	"finala/api/config"
	"fmt"
	"time"

	ms "github.com/meilisearch/meilisearch-go"
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
	Search(index string, query interface{}) (*ms.SearchResponse, error)
	CreateIndex(name string) error
	IndexExists(name string) (bool, error)
}

// meilisearchClient implements the meilisearchDescriptor interface
type meilisearchClient struct {
	client ms.ServiceManager
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
	
	// Get the primary endpoint and API key
	host := conf.Endpoints[0]
	apiKey := conf.Password
	
	// Create client using the New function (v0.32.0 style)
	client := ms.New(host, ms.WithAPIKey(apiKey))

	return &meilisearchClient{
		client: client,
		host:   host,
		apiKey: apiKey,
	}, nil
}

// healthCheck verifies the Meilisearch connection
func (m *meilisearchClient) healthCheck() error {
	// Use the client's health method directly
	_, err := m.client.Health()
	if err != nil {
		return err
	}
	
	log.Info("Successfully connected to Meilisearch")
	return nil
}

// Index implements document indexing
func (m *meilisearchClient) Index(index string, document interface{}) error {
	// Get or create the index
	idx := m.client.Index(index)
	
	// Add a single document
	_, err := idx.AddDocuments(document)
	return err
}

// Search implements search functionality
func (m *meilisearchClient) Search(index string, query interface{}) (*ms.SearchResponse, error) {
	// Convert the query params from the interface
	searchParams := query.(map[string]interface{})
	
	// Build the search query
	searchRequest := &ms.SearchRequest{
		Limit: 1000, // Default limit
	}
	
	// Extract the query string
	var q string
	if queryStr, ok := searchParams["q"].(string); ok {
		q = queryStr
	}
	
	// Add filter if present
	if filterVal, ok := searchParams["filter_by"].(string); ok && filterVal != "" {
		searchRequest.Filter = filterVal
	}
	
	// Execute the search
	idx := m.client.Index(index)
	return idx.Search(q, searchRequest)
}

// CreateIndex creates a new index in Meilisearch
func (m *meilisearchClient) CreateIndex(name string) error {
	// Create the index
	_, err := m.client.CreateIndex(&ms.IndexConfig{
		Uid:        name,
		PrimaryKey: "id",
	})
	
	if err != nil {
		return err
	}
	
	// Configure the index with filterable attributes
	return m.configureIndexSettings(name)
}

// configureIndexSettings sets up the index settings
func (m *meilisearchClient) configureIndexSettings(name string) error {
	idx := m.client.Index(name)
	const taskTimeout = 10 * time.Second // Increased timeout just in case

	// Define specific filterable attributes
	filterableAttributes := []string{
		"EventType",
		"ExecutionID",
		"ResourceName",
		// Add other known top-level or commonly filtered fields if necessary
		// For nested fields like tags, if they are under a common parent e.g., "DataMap",
		// you might list "DataMap" or rely on Meilisearch's dot notation for filtering if "*" was insufficient.
		// For now, sticking to known direct fields.
	}

	// Set filterable attributes
	taskInfoFilterable, err := idx.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		return fmt.Errorf("failed to initiate update for filterable attributes (%v): %v", filterableAttributes, err)
	}
	if _, err := m.client.WaitForTask(taskInfoFilterable.TaskUID, taskTimeout); err != nil {
		return fmt.Errorf("failed while waiting for filterable attributes update: %v", err)
	}
	log.Infof("Successfully updated filterable attributes for index %s: %v", name, filterableAttributes)

	// Set searchable attributes (keeping '*' for now, can be refined if needed)
	searchableAttributes := []string{"*"}
	taskInfoSearchable, err := idx.UpdateSearchableAttributes(&searchableAttributes)
	if err != nil {
		return fmt.Errorf("failed to initiate update for searchable attributes: %v", err)
	}
	if _, err := m.client.WaitForTask(taskInfoSearchable.TaskUID, taskTimeout); err != nil {
		return fmt.Errorf("failed while waiting for searchable attributes update: %v", err)
	}
	log.Infof("Successfully updated searchable attributes for index %s: %v", name, searchableAttributes)

	return nil
}

// IndexExists checks if an index exists
func (m *meilisearchClient) IndexExists(name string) (bool, error) {
	// Get all indexes
	indexes, err := m.client.ListIndexes(&ms.IndexesQuery{})
	if err != nil {
		return false, err
	}
	
	// Check if the index exists
	for _, index := range indexes.Results {
		if index.UID == name {
			return true, nil
		}
	}
	
	return false, nil
} 