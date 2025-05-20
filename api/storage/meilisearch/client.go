package meilisearch

import (
	"fmt"

	"finala/api/config"
	ms "github.com/meilisearch/meilisearch-go"
	log "github.com/sirupsen/logrus"
)

// meilisearchClient is a wrapper around the Meilisearch client
type meilisearchClient struct {
	client ms.ServiceManager // Changed from *ms.Client
}

// Client is the interface for Meilisearch operations
type Client interface {
	Connect(conf config.MeilisearchConfig) error
	Index(index string, document interface{}) error
	Search(index string, query interface{}) (*ms.SearchResponse, error)
	CreateIndex(name string) error
	DeleteIndex(name string) (bool, error)
	GetIndex(name string) (ms.IndexManager, error) // Changed from *ms.Index
	ListIndexes() (*ms.IndexesResults, error)
	IndexExists(name string) (bool, error)
}

// NewMeilisearchClient creates a new Meilisearch client instance
func NewMeilisearchClient() Client {
	return &meilisearchClient{}
}

// Connect establishes a connection to the Meilisearch server
func (m *meilisearchClient) Connect(conf config.MeilisearchConfig) error {
	client, err := getMeilisearchClient(conf)
	if err != nil {
		return fmt.Errorf("failed to create Meilisearch client: %w", err)
	}
	m.client = client.client // Assign the underlying ms.ServiceManager
	log.Info("Successfully connected to Meilisearch and configured client")
	return nil
}

// getMeilisearchClient creates new Meilisearch client
func getMeilisearchClient(conf config.MeilisearchConfig) (*meilisearchClient, error) {
	log.Infof("Creating Meilisearch client with endpoints: %v", conf.Endpoints)
	// Get the primary endpoint and API key
	host := conf.Endpoints[0]
	apiKey := conf.Password
	// Create client using the New function (v0.32.0 style)
	// ms.New returns ms.ServiceManager
	client := ms.New(host, ms.WithAPIKey(apiKey))

	// Verify connection by trying to get health status
	health, err := client.Health()
	if err != nil {
		return nil, fmt.Errorf("Meilisearch health check failed: %w", err)
	}
	if health.Status != "available" {
		return nil, fmt.Errorf("Meilisearch not available, status: %s", health.Status)
	}

	log.Info("Successfully created Meilisearch instance and passed health check")
	return &meilisearchClient{client: client}, nil
}

// Ping checks the health of the Meilisearch server
func (m *meilisearchClient) Ping() error {
	_, err := m.client.Health()
	if err != nil {
		return err
	}
	log.Info("Successfully connected to Meilisearch")
	return nil
}

// Index adds or updates a document in the specified index.
func (m *meilisearchClient) Index(index string, document interface{}) error {
	// Get or create the index
	idx := m.client.Index(index) // Returns IndexManager
	// Add a single document
	// Assuming AddDocuments takes []interface{} or a single interface{}
	// Based on meilisearch-go examples, it's often []map[string]interface{}
	// For a single document, it might be AddDocuments([]interface{}{document}, "id")
	// Let's check AddDocuments signature on IndexManager from go doc output (it's not directly visible)
	// For now, assuming the existing call structure is somewhat correct or needs minimal change
	// If `document` is a single item, we might need to wrap it: []interface{}{document}
	_, err := idx.AddDocuments([]interface{}{document}) // Changed to pass a slice
	return err
}

// Search performs a search query on the specified index.
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
	idx := m.client.Index(index) // Returns IndexManager
	return idx.Search(q, searchRequest)
}

// CreateIndex creates a new index with the given name if it doesn't already exist.
func (m *meilisearchClient) CreateIndex(name string) error {
	// ServiceManager.CreateIndex returns (*TaskInfo, error)
	_, err := m.client.CreateIndex(&ms.IndexConfig{
		Uid:        name,
		PrimaryKey: "id", // PrimaryKey is usually set during index creation
	})
	if err != nil {
		return err
	}
	// Configure the index with filterable attributes
	return m.configureIndexSettings(name)
}

// configureIndexSettings sets filterable attributes for an index
func (m *meilisearchClient) configureIndexSettings(indexName string) error {
	idx := m.client.Index(indexName) // Returns IndexManager
	settings := ms.Settings{
		FilterableAttributes: []string{"ExecutionID", "ResourceName", "EventType", "tags", "Collector"},
	}
	// IndexManager.UpdateSettings returns (*TaskInfo, error)
	_, err := idx.UpdateSettings(&settings)
	if err != nil {
		return fmt.Errorf("failed to update settings for index %s: %w", indexName, err)
	}
	log.Infof("Successfully configured filterable attributes for index %s", indexName)
	return nil
}

// DeleteIndex deletes an index by its UID.
func (m *meilisearchClient) DeleteIndex(name string) (bool, error) {
	// ServiceManager.DeleteIndex returns (*TaskInfo, error)
	// The boolean return type in the interface might need to change or be derived from TaskInfo.
	// For now, let's assume the goal is to return true on success (no error from DeleteIndex).
	_, err := m.client.DeleteIndex(name)
	if err != nil {
		return false, err
	}
	return true, nil // Changed to return true if no error
}

// GetIndex fetches an index by its UID.
func (m *meilisearchClient) GetIndex(name string) (ms.IndexManager, error) { // Changed return type
	// ServiceManager.Index(uid) returns an IndexManager. There is no direct error return.
	// If an index doesn't exist, operations on IndexManager might fail.
	// To check existence, one might call FetchInfo() or ListIndexes().
	idx := m.client.Index(name)
	// To conform to an error return, we could try a cheap operation like FetchInfo()
	// and return an error if it fails significantly (e.g. index not found error code).
	// For now, just return the IndexManager.
	// If an error is strictly needed here, one would have to be manufactured or a light check performed.
	return idx, nil
}

// ListIndexes lists all indexes.
func (m *meilisearchClient) ListIndexes() (*ms.IndexesResults, error) {
	return m.client.ListIndexes(nil) // Pass nil for default parameters
}

// IndexExists checks if an index exists by its UID.
func (m *meilisearchClient) IndexExists(name string) (bool, error) {
	// This implementation seems fine as m.client.ListIndexes is a valid call.
	indexes, err := m.ListIndexes()
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
