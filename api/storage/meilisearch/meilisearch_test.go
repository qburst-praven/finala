package meilisearch

import (
	"errors"
	"testing"
	"time"

	"finala/api/config"
	ms "github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewStorageManager_Success is skipped. Directly testing NewStorageManager is complex
// due to its internal client creation and the Connect call.
// Refactoring NewStorageManager to accept a Client interface (dependency injection)
// would make this more unit-testable. Otherwise, this is better for integration tests.
func TestNewStorageManager_Success(t *testing.T) {
	t.Skip("Skipping test for NewStorageManager success; requires integration setup or refactor for DI.")
}

// TestNewStorageManager_ConnectFails is also skipped for similar reasons.
func TestNewStorageManager_ConnectFails(t *testing.T) {
	t.Skip("Skipping test for NewStorageManager connection failure; requires integration setup or refactor for DI.")
}

// TestStorageManager_setCreateCurrentIndexDay_IndexExists tests when the daily index already exists.
func TestStorageManager_setCreateCurrentIndexDay_IndexExists(t *testing.T) {
	mockClient := new(MockClient) // Using the local MockClient for StorageManager.client
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"}, // Base name
	}
	expectedIndexName := "events_" + time.Now().Format("2006_01_02")

	mockClient.On("IndexExists", expectedIndexName).Return(true, nil).Once()

	success := sm.setCreateCurrentIndexDay()
	assert.True(t, success, "setCreateCurrentIndexDay should succeed if index exists")
	assert.Equal(t, expectedIndexName, sm.currentIndexDay, "currentIndexDay should be set to the existing index name")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_setCreateCurrentIndexDay_CreatesNewIndex tests when the daily index does not exist and is created.
func TestStorageManager_setCreateCurrentIndexDay_CreatesNewIndex(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}
	expectedIndexName := "events_" + time.Now().Format("2006_01_02")

	mockClient.On("IndexExists", expectedIndexName).Return(false, nil).Once()
	mockClient.On("CreateIndex", expectedIndexName).Return(nil).Once()

	success := sm.setCreateCurrentIndexDay()
	assert.True(t, success, "setCreateCurrentIndexDay should succeed if new index is created")
	assert.Equal(t, expectedIndexName, sm.currentIndexDay, "currentIndexDay should be set to the new index name")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_setCreateCurrentIndexDay_IndexExistsFails tests failure when checking if index exists.
func TestStorageManager_setCreateCurrentIndexDay_IndexExistsFails(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}
	expectedIndexName := "events_" + time.Now().Format("2006_01_02")
	expectedError := errors.New("failed to check index existence")

	mockClient.On("IndexExists", expectedIndexName).Return(false, expectedError).Once()

	success := sm.setCreateCurrentIndexDay()
	assert.False(t, success, "setCreateCurrentIndexDay should fail if IndexExists fails")
	assert.Empty(t, sm.currentIndexDay, "currentIndexDay should be empty on failure")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_setCreateCurrentIndexDay_CreateIndexFails tests failure when creating a new index.
func TestStorageManager_setCreateCurrentIndexDay_CreateIndexFails(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}
	expectedIndexName := "events_" + time.Now().Format("2006_01_02")
	expectedError := errors.New("failed to create index")

	mockClient.On("IndexExists", expectedIndexName).Return(false, nil).Once()
	mockClient.On("CreateIndex", expectedIndexName).Return(expectedError).Once()

	success := sm.setCreateCurrentIndexDay()
	assert.False(t, success, "setCreateCurrentIndexDay should fail if CreateIndex fails")
	assert.Empty(t, sm.currentIndexDay, "currentIndexDay should be empty on failure")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_AddEvent_Success tests successful event addition.
func TestStorageManager_AddEvent_Success(t *testing.T) {
	mockClient := new(MockClient)
	currentIndex := "events_2023_01_01"
	sm := &StorageManager{
		client:          mockClient,
		currentIndexDay: currentIndex, // Pre-set for the test
	}
	event := map[string]interface{}{"id": "1", "data": "test data"}

	mockClient.On("Index", currentIndex, event).Return(nil).Once()

	err := sm.AddEvent(event)
	assert.NoError(t, err, "AddEvent should not return an error on success")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_AddEvent_Failure tests event addition failure.
func TestStorageManager_AddEvent_Failure(t *testing.T) {
	mockClient := new(MockClient)
	currentIndex := "events_2023_01_01"
	sm := &StorageManager{
		client:          mockClient,
		currentIndexDay: currentIndex,
	}
	event := map[string]interface{}{"id": "1", "data": "test data"}
	expectedError := errors.New("failed to index event")

	mockClient.On("Index", currentIndex, event).Return(expectedError).Once()

	err := sm.AddEvent(event)
	assert.Error(t, err, "AddEvent should return an error on failure")
	assert.Equal(t, expectedError, err, "Error should match the one from underlying Index call")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_AddEvent_NoCurrentIndex tests adding an event when currentIndexDay is not set.
func TestStorageManager_AddEvent_NoCurrentIndex(t *testing.T) {
	mockClient := new(MockClient) // MockClient is not actually used here as it should fail before
	sm := &StorageManager{
		client: mockClient,
		// currentIndexDay is deliberately not set
	}
	event := map[string]interface{}{"id": "1", "data": "test data"}

	err := sm.AddEvent(event)
	assert.Error(t, err, "AddEvent should return an error if currentIndexDay is not set")
	assert.Contains(t, err.Error(), "current index for today is not set", "Error message should indicate missing index")
	// No expectations on mockClient as it shouldn't be called.
}

// TestStorageManager_SearchEvents_Success tests successful event search.
func TestStorageManager_SearchEvents_Success(t *testing.T) {
	mockClient := new(MockClient)
	currentIndex := "events_search_index"
	sm := &StorageManager{
		client:          mockClient,
		currentIndexDay: currentIndex, // Assuming search targets current day's index for simplicity
		config:          config.StorageMeilisearch{EventsIndexName: "events"}, // Needed for GetSearchIndexNames
	}
	query := map[string]interface{}{"q": "test query"}
	expectedResponse := &ms.SearchResponse{Hits: []interface{}{map[string]interface{}{"id": "res1"}}}

	// GetSearchIndexNames will be called. Let's assume it returns the current index for this test.
	// A more robust test for GetSearchIndexNames would be separate.
	mockClient.On("ListIndexes").Return(&ms.IndexesResults{Results: []ms.Index{{UID: currentIndex}}}, nil).Maybe() // Allow flexibility
	mockClient.On("Search", currentIndex, query).Return(expectedResponse, nil).Once()

	resp, err := sm.SearchEvents(query)
	assert.NoError(t, err, "SearchEvents should not return an error on success")
	assert.Equal(t, expectedResponse, resp, "SearchEvents response should match expected")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_SearchEvents_Failure tests event search failure.
func TestStorageManager_SearchEvents_Failure(t *testing.T) {
	mockClient := new(MockClient)
	currentIndex := "events_search_index_fail"
	sm := &StorageManager{
		client:          mockClient,
		currentIndexDay: currentIndex,
		config:          config.StorageMeilisearch{EventsIndexName: "events"},
	}
	query := map[string]interface{}{"q": "find this"}
	expectedError := errors.New("search operation failed")

	mockClient.On("ListIndexes").Return(&ms.IndexesResults{Results: []ms.Index{{UID: currentIndex}}}, nil).Maybe()
	mockClient.On("Search", currentIndex, query).Return(nil, expectedError).Once()

	resp, err := sm.SearchEvents(query)
	assert.Error(t, err, "SearchEvents should return an error on failure")
	assert.Nil(t, resp, "SearchEvents response should be nil on failure")
	assert.Equal(t, expectedError, err, "Error should match the one from underlying Search call")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_SearchEvents_NoSearchableIndexes tests search when no relevant indexes are found.
func TestStorageManager_SearchEvents_NoSearchableIndexes(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}
	query := map[string]interface{}{"q": "anything"}

	// Simulate GetSearchIndexNames returning an empty list
	mockClient.On("ListIndexes").Return(&ms.IndexesResults{Results: []ms.Index{}}, nil).Once()

	resp, err := sm.SearchEvents(query)
	assert.NoError(t, err, "SearchEvents should not error if no indexes found, but return empty response")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.Empty(t, resp.Hits, "Response Hits should be empty if no indexes to search")
	mockClient.AssertExpectations(t)
}

// TestStorageManager_GetSearchIndexNames_NoFilter tests GetSearchIndexNames with no date filter.
func TestStorageManager_GetSearchIndexNames_NoFilter(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}

	allIndexes := []ms.Index{
		{UID: "events_2023_01_01"},
		{UID: "events_2023_01_02"},
		{UID: "other_index"},
	}
	expectedSearchableIndexes := []string{"events_2023_01_01", "events_2023_01_02"}

	mockClient.On("ListIndexes").Return(&ms.IndexesResults{Results: allIndexes}, nil).Once()

	searchableIndexes, err := sm.GetSearchIndexNames(nil, nil) // No date filters
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSearchableIndexes, searchableIndexes)
	mockClient.AssertExpectations(t)
}

// TestStorageManager_GetSearchIndexNames_WithDateFilters tests GetSearchIndexNames with date filters.
func TestStorageManager_GetSearchIndexNames_WithDateFilters(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}

	fromDate := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)

	allIndexes := []ms.Index{
		{UID: "events_2023_01_01"}, // Before fromDate
		{UID: "events_2023_01_02"}, // Matches fromDate
		{UID: "events_2023_01_03"}, // Matches toDate (exclusive, so not included if exact match for 'days ago')
		{UID: "events_2023_01_04"}, // After toDate
		{UID: "non_events_index_2023_01_02"},
		{UID: "events_invalid_date"},
	}
	// For GetSearchIndexNames, toDate is exclusive for the day part for simplicity of daysAgo logic.
	// If we filter for 2023-01-02 to 2023-01-03, we expect events_2023_01_02.
	expectedSearchableIndexes := []string{"events_2023_01_02"}

	mockClient.On("ListIndexes").Return(&ms.IndexesResults{Results: allIndexes}, nil).Once()

	searchableIndexes, err := sm.GetSearchIndexNames(&fromDate, &toDate)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSearchableIndexes, searchableIndexes)
	mockClient.AssertExpectations(t)
}

// TestStorageManager_GetSearchIndexNames_ListIndexesFails tests GetSearchIndexNames when ListIndexes fails.
func TestStorageManager_GetSearchIndexNames_ListIndexesFails(t *testing.T) {
	mockClient := new(MockClient)
	sm := &StorageManager{
		client: mockClient,
		config: config.StorageMeilisearch{EventsIndexName: "events"},
	}
	expectedError := errors.New("failed to list indexes")

	mockClient.On("ListIndexes").Return(nil, expectedError).Once()

	searchableIndexes, err := sm.GetSearchIndexNames(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, searchableIndexes)
	assert.Equal(t, expectedError, err)
	mockClient.AssertExpectations(t)
}

// TestGetIndexNameByDate verifies the utility function for generating index names.
func TestGetIndexNameByDate(t *testing.T) {
	baseName := "my_events"
	date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	expectedName := "my_events_2024_03_15"
	assert.Equal(t, expectedName, GetIndexNameByDate(baseName, date))
}

// TestGetTimeFromIndexDateString_Valid tests parsing a valid date string from an index name.
func TestGetTimeFromIndexDateString_Valid(t *testing.T) {
	dateStr := "2024_03_15"
	expectedTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	t, err := getTimeFromIndexDateString(dateStr)
	assert.NoError(t, err)
	assert.Equal(t, expectedTime, t)
}

// TestGetTimeFromIndexDateString_Invalid tests parsing an invalid date string.
func TestGetTimeFromIndexDateString_Invalid(t *testing.T) {
	dateStr := "invalid_date_format"
	_, err := getTimeFromIndexDateString(dateStr)
	assert.Error(t, err, "Should return an error for invalid date format")
} 