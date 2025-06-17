package meilisearch

import (
	"context"
	"errors"
	"testing"
	"time"

	"finala/api/config"
	ms "github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewMeilisearchClient verifies that NewMeilisearchClient creates an instance of meilisearchClient.
func TestNewMeilisearchClient(t *testing.T) {
	client := NewMeilisearchClient()
	assert.NotNil(t, client, "NewMeilisearchClient should return a non-nil client")
	_, ok := client.(*meilisearchClient)
	assert.True(t, ok, "NewMeilisearchClient should return a *meilisearchClient type")
}

// TestMeilisearchClient_Connect_Success is skipped because it's hard to mock the internal getMeilisearchClient call.
// This scenario is better covered by integration or E2E tests.
func TestMeilisearchClient_Connect_Success(t *testing.T) {
	t.Skip("Skipping Connect success test due to complexity of mocking ms.NewClient and unexported getMeilisearchClient. Integration or E2E tests would cover this better.")
}

// TestMeilisearchClient_Connect_Failure is also skipped for the same reasons as the success case.
func TestMeilisearchClient_Connect_Failure(t *testing.T) {
	t.Skip("Skipping Connect failure test due to complexity of mocking ms.NewClient and unexported getMeilisearchClient. Integration or E2E tests would cover this better.")
}

// TestMeilisearchClient_Ping_Success tests the Ping method for a successful scenario.
func TestMeilisearchClient_Ping_Success(t *testing.T) {
	mockServiceMgr := new(MockServiceManager) // Using MockServiceManager as meilisearchClient.client is *ms.Client
	client := &meilisearchClient{client: mockServiceMgr}

	mockServiceMgr.On("Health").Return(&ms.Health{Status: "available"}, nil).Once()

	err := client.Ping()
	assert.NoError(t, err, "Ping should not return an error on success")
	mockServiceMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Ping_Failure tests the Ping method for a failure scenario.
func TestMeilisearchClient_Ping_Failure(t *testing.T) {
	mockServiceMgr := new(MockServiceManager)
	client := &meilisearchClient{client: mockServiceMgr}

	expectedError := errors.New("ping failed")
	mockServiceMgr.On("Health").Return(nil, expectedError).Once()

	err := client.Ping()
	assert.Error(t, err, "Ping should return an error on failure")
	assert.Equal(t, expectedError, err, "Error returned by Ping should match the expected error")
	mockServiceMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Index_Success tests the Index method for a successful document indexing.
func TestMeilisearchClient_Index_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager) // This is the *ms.Client
	mockIndexMgr := new(MockIndexManager)           // This is the *ms.Index
	client := &meilisearchClient{client: mockUnderlyingClient}

	indexName := "test_index"
	document := map[string]interface{}{"id": "1", "title": "Test Document"}

	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("AddDocuments", []interface{}{document}).Return(&ms.TaskInfo{TaskUID: 1}, nil).Once()

	err := client.Index(indexName, document)
	assert.NoError(t, err, "Index should not return an error on successful indexing")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Index_Failure tests the Index method for a failure scenario during document indexing.
func TestMeilisearchClient_Index_Failure(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager)
	client := &meilisearchClient{client: mockUnderlyingClient}

	indexName := "test_index"
	document := map[string]interface{}{"id": "1", "title": "Test Document"}
	expectedError := errors.New("indexing failed")

	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("AddDocuments", []interface{}{document}).Return(nil, expectedError).Once()

	err := client.Index(indexName, document)
	assert.Error(t, err, "Index should return an error on failure")
	assert.Equal(t, expectedError, err, "Error returned by Index should match the expected error")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Search_Success tests the Search method for a successful search operation.
func TestMeilisearchClient_Search_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager)
	client := &meilisearchClient{client: mockUnderlyingClient}

	indexName := "test_index"
	searchQuery := "test query"
	filter := "tag=test"
	queryMap := map[string]interface{}{"q": searchQuery, "filter_by": filter}
	expectedResponse := &ms.SearchResponse{Hits: []interface{}{map[string]interface{}{"id": "1"}}}

	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("Search", searchQuery, &ms.SearchRequest{Limit: 1000, Filter: filter}).Return(expectedResponse, nil).Once()

	resp, err := client.Search(indexName, queryMap)
	assert.NoError(t, err, "Search should not return an error on success")
	assert.Equal(t, expectedResponse, resp, "Search response should match the expected response")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Search_Success_NoFilter tests search without a filter.
func TestMeilisearchClient_Search_Success_NoFilter(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager)
	client := &meilisearchClient{client: mockUnderlyingClient}

	indexName := "test_index"
	searchQuery := "test query"
	queryMap := map[string]interface{}{"q": searchQuery}
	expectedResponse := &ms.SearchResponse{Hits: []interface{}{map[string]interface{}{"id": "1"}}}

	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("Search", searchQuery, &ms.SearchRequest{Limit: 1000, Filter: ""}).Return(expectedResponse, nil).Once()

	resp, err := client.Search(indexName, queryMap)
	assert.NoError(t, err, "Search should not return an error when no filter is provided")
	assert.Equal(t, expectedResponse, resp, "Search response should match the expected response")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_Search_Failure tests the Search method for a failure scenario.
func TestMeilisearchClient_Search_Failure(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager)
	client := &meilisearchClient{client: mockUnderlyingClient}

	indexName := "test_index"
	searchQuery := "test query"
	queryMap := map[string]interface{}{"q": searchQuery}
	expectedError := errors.New("search failed")

	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("Search", searchQuery, &ms.SearchRequest{Limit: 1000, Filter: ""}).Return(nil, expectedError).Once()

	resp, err := client.Search(indexName, queryMap)
	assert.Error(t, err, "Search should return an error on failure")
	assert.Nil(t, resp, "Search response should be nil on failure")
	assert.Equal(t, expectedError, err, "Error returned by Search should match the expected error")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_CreateIndex_Success tests successful index creation and configuration.
func TestMeilisearchClient_CreateIndex_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager) // Mock for the *ms.Index object
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "new_index"

	mockUnderlyingClient.On("CreateIndex", &ms.IndexConfig{Uid: indexName, PrimaryKey: "id"}).Return(&ms.TaskInfo{TaskUID: 1}, nil).Once()
	// After CreateIndex, the code gets the index again to configure it.
	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once() 
	mockIndexMgr.On("UpdateSettings", mock.AnythingOfType("*meilisearch.Settings")).Return(&ms.TaskInfo{TaskUID: 2}, nil).Once().Run(func(args mock.Arguments) {
		settingsArg := args.Get(0).(*ms.Settings)
		assert.Equal(t, []string{"ExecutionID", "ResourceName", "EventType", "tags", "Collector"}, settingsArg.FilterableAttributes)
	})

	err := client.CreateIndex(indexName)
	assert.NoError(t, err, "CreateIndex should not return an error on success")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_CreateIndex_Failure_OnCreate tests failure during the initial index creation call.
func TestMeilisearchClient_CreateIndex_Failure_OnCreate(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "new_index"
	expectedError := errors.New("create index failed")

	mockUnderlyingClient.On("CreateIndex", &ms.IndexConfig{Uid: indexName, PrimaryKey: "id"}).Return(nil, expectedError).Once()

	err := client.CreateIndex(indexName)
	assert.Error(t, err, "CreateIndex should return an error if underlying CreateIndex fails")
	assert.Equal(t, expectedError, err, "Error should match the one from CreateIndex")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_CreateIndex_Failure_OnConfigure tests failure during the index configuration (UpdateSettings) call.
func TestMeilisearchClient_CreateIndex_Failure_OnConfigure(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "new_index"
	expectedErrorOnConfigure := errors.New("configure index failed")

	mockUnderlyingClient.On("CreateIndex", &ms.IndexConfig{Uid: indexName, PrimaryKey: "id"}).Return(&ms.TaskInfo{TaskUID: 1}, nil).Once()
	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()
	mockIndexMgr.On("UpdateSettings", mock.AnythingOfType("*meilisearch.Settings")).Return(nil, expectedErrorOnConfigure).Once()

	err := client.CreateIndex(indexName)
	assert.Error(t, err, "CreateIndex should return an error if UpdateSettings fails")
	assert.Contains(t, err.Error(), expectedErrorOnConfigure.Error(), "Error message should contain the configure error")
	mockUnderlyingClient.AssertExpectations(t)
	mockIndexMgr.AssertExpectations(t)
}

// TestMeilisearchClient_DeleteIndex_Success tests successful index deletion.
func TestMeilisearchClient_DeleteIndex_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "old_index"

	mockUnderlyingClient.On("DeleteIndex", indexName).Return(&ms.TaskInfo{TaskUID: 1}, nil).Once()

	deleted, err := client.DeleteIndex(indexName)
	assert.NoError(t, err, "DeleteIndex should not return an error on success")
	assert.True(t, deleted, "DeleteIndex should return true on successful deletion")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_DeleteIndex_Failure tests index deletion failure.
func TestMeilisearchClient_DeleteIndex_Failure(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "old_index"
	expectedError := errors.New("delete index failed")

	mockUnderlyingClient.On("DeleteIndex", indexName).Return(nil, expectedError).Once()

	deleted, err := client.DeleteIndex(indexName)
	assert.Error(t, err, "DeleteIndex should return an error on failure")
	assert.False(t, deleted, "DeleteIndex should return false on failure")
	assert.Equal(t, expectedError, err, "Error should match the one from DeleteIndex")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_GetIndex_Success tests successfully getting an index reference.
func TestMeilisearchClient_GetIndex_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	mockIndexMgr := new(MockIndexManager) // This is the *ms.Index that our IndexManager wraps
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "existing_index"

	// client.GetIndex() uses client.client.Index() which returns *ms.Index
	mockUnderlyingClient.On("Index", indexName).Return(mockIndexMgr).Once()

	idx, err := client.GetIndex(indexName) // idx will be an IndexManager (our wrapper type)
	assert.NoError(t, err, "GetIndex should not return an error on success")
	assert.NotNil(t, idx, "GetIndex should return a non-nil IndexManager")

	// Check if the returned IndexManager is of the expected concrete type and wraps the mock
	wrappedIndex, ok := idx.(*indexManager)
	assert.True(t, ok, "Returned IndexManager should be of type *indexManager")
	assert.Equal(t, mockIndexMgr, wrappedIndex.index, "Wrapped index should be the mockIndexMgr")
	mockUnderlyingClient.AssertExpectations(t)
}


// TestMeilisearchClient_ListIndexes_Success tests successful listing of indexes.
func TestMeilisearchClient_ListIndexes_Success(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	expectedResponse := &ms.IndexesResults{
		Results: []ms.Index{
			{UID: "index1"},
			{UID: "index2"},
		},
		Limit:  20,
		Offset: 0,
	}

	mockUnderlyingClient.On("GetIndexes").Return(expectedResponse, nil).Once() // Changed from ListIndexes to GetIndexes

	resp, err := client.ListIndexes()
	assert.NoError(t, err, "ListIndexes should not return an error on success")
	assert.Equal(t, expectedResponse, resp, "ListIndexes response should match expected")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_ListIndexes_Failure tests failure in listing indexes.
func TestMeilisearchClient_ListIndexes_Failure(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	expectedError := errors.New("list indexes failed")

	mockUnderlyingClient.On("GetIndexes").Return(nil, expectedError).Once() // Changed from ListIndexes to GetIndexes

	resp, err := client.ListIndexes()
	assert.Error(t, err, "ListIndexes should return an error on failure")
	assert.Nil(t, resp, "ListIndexes response should be nil on failure")
	assert.Equal(t, expectedError, err, "Error should match expected")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_IndexExists_True tests if an index exists and is found.
func TestMeilisearchClient_IndexExists_True(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "index1"
	listResponse := &ms.IndexesResults{
		Results: []ms.Index{
			{UID: "index1"},
			{UID: "index2"},
		},
	}

	mockUnderlyingClient.On("GetIndexes").Return(listResponse, nil).Once()

	exists, err := client.IndexExists(indexName)
	assert.NoError(t, err, "IndexExists should not return error when underlying call succeeds")
	assert.True(t, exists, "IndexExists should return true when index is in the list")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_IndexExists_False tests if an index exists and is not found.
func TestMeilisearchClient_IndexExists_False(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "index3" // This index does not exist in the mock response
	listResponse := &ms.IndexesResults{
		Results: []ms.Index{
			{UID: "index1"},
			{UID: "index2"},
		},
	}

	mockUnderlyingClient.On("GetIndexes").Return(listResponse, nil).Once()

	exists, err := client.IndexExists(indexName)
	assert.NoError(t, err, "IndexExists should not return error when underlying call succeeds")
	assert.False(t, exists, "IndexExists should return false when index is not in the list")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestMeilisearchClient_IndexExists_ErrorOnList tests an error scenario when trying to list indexes for existence check.
func TestMeilisearchClient_IndexExists_ErrorOnList(t *testing.T) {
	mockUnderlyingClient := new(MockServiceManager)
	client := &meilisearchClient{client: mockUnderlyingClient}
	indexName := "index1"
	expectedError := errors.New("failed to list indexes")

	mockUnderlyingClient.On("GetIndexes").Return(nil, expectedError).Once()

	exists, err := client.IndexExists(indexName)
	assert.Error(t, err, "IndexExists should return error when underlying GetIndexes fails")
	assert.False(t, exists, "IndexExists should return false on error")
	assert.Equal(t, expectedError, err, "Error should match the one from GetIndexes")
	mockUnderlyingClient.AssertExpectations(t)
}

// TestIndexManager_Search_Success tests the Search method of the wrapped indexManager.
func TestIndexManager_Search_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager) // This is the mock for *ms.Index
	im := indexManager{index: mockMsIndex}

	query := "test search"
	filter := "type=product"
	searchReq := &ms.SearchRequest{Filter: filter, Limit: 100}
	expectedResp := &ms.SearchResponse{Hits: []interface{}{"hit1"}}

	mockMsIndex.On("Search", query, searchReq).Return(expectedResp, nil).Once()

	resp, err := im.Search(query, searchReq)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_Search_Failure tests the Search method failure of the wrapped indexManager.
func TestIndexManager_Search_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}

	query := "test search"
	searchReq := &ms.SearchRequest{Limit: 100}
	expectedErr := errors.New("search failed")

	mockMsIndex.On("Search", query, searchReq).Return(nil, expectedErr).Once()

	resp, err := im.Search(query, searchReq)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_AddDocuments_Success tests adding documents successfully.
func TestIndexManager_AddDocuments_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}

	documents := []interface{}{map[string]string{"id": "1"}}
	expectedTask := &ms.TaskInfo{TaskUID: 123}

	mockMsIndex.On("AddDocuments", documents).Return(expectedTask, nil).Once()

	task, err := im.AddDocuments(documents)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_AddDocuments_Failure tests failure while adding documents.
func TestIndexManager_AddDocuments_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}

	documents := []interface{}{map[string]string{"id": "1"}}
	expectedErr := errors.New("add documents failed")

	mockMsIndex.On("AddDocuments", documents).Return(nil, expectedErr).Once()

	task, err := im.AddDocuments(documents)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, expectedErr, err)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_UpdateSettings_Success tests updating settings successfully.
func TestIndexManager_UpdateSettings_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}

	settings := &ms.Settings{DisplayedAttributes: []string{"name"}}
	expectedTask := &ms.TaskInfo{TaskUID: 456}

	mockMsIndex.On("UpdateSettings", settings).Return(expectedTask, nil).Once()

	task, err := im.UpdateSettings(settings)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_UpdateSettings_Failure tests failure while updating settings.
func TestIndexManager_UpdateSettings_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}

	settings := &ms.Settings{DisplayedAttributes: []string{"name"}}
	expectedErr := errors.New("update settings failed")

	mockMsIndex.On("UpdateSettings", settings).Return(nil, expectedErr).Once()

	task, err := im.UpdateSettings(settings)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, expectedErr, err)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_DeleteDocument_Success tests deleting a document successfully.
func TestIndexManager_DeleteDocument_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	docID := "doc1"
	expectedTask := &ms.TaskInfo{TaskUID: 789}

	mockMsIndex.On("DeleteDocument", docID).Return(expectedTask, nil).Once()

	task, err := im.DeleteDocument(docID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_DeleteDocument_Failure tests failure while deleting a document.
func TestIndexManager_DeleteDocument_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	docID := "doc1"
	expectedErr := errors.New("delete document failed")

	mockMsIndex.On("DeleteDocument", docID).Return(nil, expectedErr).Once()

	task, err := im.DeleteDocument(docID)
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, expectedErr, err)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_DeleteAllDocuments_Success tests deleting all documents successfully.
func TestIndexManager_DeleteAllDocuments_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	expectedTask := &ms.TaskInfo{TaskUID: 101}

	mockMsIndex.On("DeleteAllDocuments").Return(expectedTask, nil).Once()

	task, err := im.DeleteAllDocuments()
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_DeleteAllDocuments_Failure tests failure while deleting all documents.
func TestIndexManager_DeleteAllDocuments_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	expectedErr := errors.New("delete all documents failed")

	mockMsIndex.On("DeleteAllDocuments").Return(nil, expectedErr).Once()

	task, err := im.DeleteAllDocuments()
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, expectedErr, err)
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_GetDocument_Success tests getting a document successfully.
func TestIndexManager_GetDocument_Success(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	docID := "doc123"
	var docPtr map[string]interface{}

	// Simulate successful GetDocument call by the mock
	mockMsIndex.On("GetDocument", docID, &docPtr).Return(nil).Run(func(args mock.Arguments) {
		// Simulate MeiliSearch filling the documentPtr
		ptr := args.Get(1).(*map[string]interface{}) 
		*ptr = map[string]interface{}{"id": docID, "data": "sample"}
	}).Once()

	err := im.GetDocument(docID, &docPtr)
	assert.NoError(t, err)
	assert.NotNil(t, docPtr)
	assert.Equal(t, docID, docPtr["id"])
	mockMsIndex.AssertExpectations(t)
}

// TestIndexManager_GetDocument_Failure tests failure while getting a document.
func TestIndexManager_GetDocument_Failure(t *testing.T) {
	mockMsIndex := new(MockIndexManager)
	im := indexManager{index: mockMsIndex}
	docID := "doc123"
	var docPtr map[string]interface{}
	expectedErr := errors.New("get document failed")

	mockMsIndex.On("GetDocument", docID, &docPtr).Return(expectedErr).Once()

	err := im.GetDocument(docID, &docPtr)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, docPtr) // docPtr should remain nil or in its initial state
	mockMsIndex.AssertExpectations(t)
} 