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
func NewClient(conf config.ElasticsearchConfig) (*meilisearchClient, error) {
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
func getMeilisearchClient(conf config.ElasticsearchConfig) (*meilisearchClient, error) {
	log.Infof("Creating Meilisearch client with endpoints: %v", conf.Endpoints)
	
	// For v0.12.0, set the apiKey directly and not through the config
	// since many methods might not be using the config correctly
	host := conf.Endpoints[0]
	apiKey := conf.Password
	
	// Create client with empty config, we'll manually add Authorization headers
	config := meilisearch.Config{
		Host: host,
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
		return fmt.Errorf("health check failed with status: %s", resp.Status)
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
		return nil, err
	}
	defer resp.Body.Close()
	
	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Search failed: status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	// Parse the response
	var searchResult meilisearch.SearchResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return nil, err
	}
	
	return &searchResult, nil
}

// CreateIndex creates a new index
func (m *meilisearchClient) CreateIndex(name string) error {
	// Check if index exists first
	exists, err := m.IndexExists(name)
	if err != nil {
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
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Failed to create index: %s, status: %d, body: %s", name, resp.StatusCode, string(body))
	}
	
	log.Infof("Index %s created successfully", name)
	return m.configureIndexSettings(name)
}

// configureIndexSettings sets up the necessary index settings for optimal search
func (m *meilisearchClient) configureIndexSettings(name string) error {
	// Configure searchable attributes using direct HTTP request
	searchableAttrs := []string{
		"ResourceName",
		"ExecutionID",
		"EventType",
		"Data",
		"*",  // Make all fields searchable
	}
	
	if err := m.updateSettings(name, "searchableAttributes", searchableAttrs); err != nil {
		return err
	}
	
	// Configure ranking rules
	rankingRules := []string{
		"typo",
		"words",
		"proximity",
		"attribute",
		"exactness",
		"desc(Timestamp)",
	}
	
	if err := m.updateSettings(name, "rankingRules", rankingRules); err != nil {
		return err
	}
	
	// Configure attributes for faceting
	facetingAttrs := []string{
		"ResourceName",
		"ExecutionID", 
		"EventType",
		"Timestamp",
		"id",
	}
	
	if err := m.updateSettings(name, "attributesForFaceting", facetingAttrs); err != nil {
		return err
	}
	
	log.Infof("Successfully configured index settings for %s", name)
	return nil
}

// updateSettings is a helper method to update various settings
func (m *meilisearchClient) updateSettings(indexName, settingName string, value interface{}) error {
	url := fmt.Sprintf("%s/indexes/%s/settings/%s", m.host, indexName, settingName)
	
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonValue))
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
		return fmt.Errorf("Failed to update %s: status: %d, body: %s", settingName, resp.StatusCode, string(body))
	}
	
	return nil
}

// IndexExists checks if an index exists
func (m *meilisearchClient) IndexExists(name string) (bool, error) {
	// For v0.12.0, make a direct HTTP request to check if index exists
	url := fmt.Sprintf("%s/indexes/%s", m.host, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	
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
	return false, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
} 