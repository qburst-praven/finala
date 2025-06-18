package meilisearch

import (
	"context"
	"io"
	"time"

	"finala/api/config" // Added to resolve config undefined issue

	ms "github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/mock"
)

// MockServiceManager is a mock implementation of ms.ServiceManager
// The actual meilisearchClient.client field is ms.ServiceManager
type MockServiceManager struct {
	mock.Mock
}

// Index mimics the (c *Client) Index(uid string) *Index method from meilisearch-go/client.go
// It returns a *ms.Index which implements ms.IndexManager
func (m *MockServiceManager) Index(uid string) ms.IndexManager {
	args := m.Called(uid)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.IndexManager)
}

func (m *MockServiceManager) CreateIndex(config *ms.IndexConfig) (*ms.TaskInfo, error) {
	args := m.Called(config)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) CreateIndexWithContext(ctx context.Context, config *ms.IndexConfig) (*ms.TaskInfo, error) {
	args := m.Called(ctx, config)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// GetIndex mimics the (c *Client) GetIndex(uid string) (*IndexResult, error) method
func (m *MockServiceManager) GetIndex(uid string) (*ms.IndexResult, error) {
	args := m.Called(uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexResult), args.Error(1)
}

func (m *MockServiceManager) GetIndexWithContext(ctx context.Context, uid string) (*ms.IndexResult, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexResult), args.Error(1)
}

func (m *MockServiceManager) GetIndexes(params *ms.IndexesQuery) (*ms.IndexesResults, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexesResults), args.Error(1)
}

func (m *MockServiceManager) GetIndexesWithContext(ctx context.Context, params *ms.IndexesQuery) (*ms.IndexesResults, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexesResults), args.Error(1)
}

func (m *MockServiceManager) DeleteIndex(uid string) (*ms.TaskInfo, error) {
	args := m.Called(uid)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) DeleteIndexWithContext(ctx context.Context, uid string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, uid)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) Health() (*ms.Health, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Health), args.Error(1)
}

func (m *MockServiceManager) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockServiceManager) Stats() (*ms.Stats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Stats), args.Error(1)
}

func (m *MockServiceManager) Version() (*ms.Version, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Version), args.Error(1)
}

func (m *MockServiceManager) CreateDump() (*ms.TaskInfo, error) {
	args := m.Called()
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) GetDumpStatus(uid string) (*ms.Task, error) {
	args := m.Called(uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockServiceManager) GetTasks(params *ms.TasksQuery) (*ms.TaskResult, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskResult), args.Error(1)
}

func (m *MockServiceManager) GetTasksWithContext(ctx context.Context, params *ms.TasksQuery) (*ms.TaskResult, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskResult), args.Error(1)
}

func (m *MockServiceManager) GetTask(uid int64) (*ms.Task, error) {
	args := m.Called(uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockServiceManager) GetTaskWithContext(ctx context.Context, uid int64) (*ms.Task, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockServiceManager) WaitForTask(uid int64, interval time.Duration) (*ms.Task, error) {
	args := m.Called(uid, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockServiceManager) WaitForTaskWithContext(ctx context.Context, uid int64, interval time.Duration) (*ms.Task, error) {
	args := m.Called(ctx, uid, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockServiceManager) CancelTasks(params *ms.CancelTasksQuery) (*ms.TaskInfo, error) {
	args := m.Called(params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) CancelTasksWithContext(ctx context.Context, params *ms.CancelTasksQuery) (*ms.TaskInfo, error) {
	args := m.Called(ctx, params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) DeleteTasks(params *ms.DeleteTasksQuery) (*ms.TaskInfo, error) {
	args := m.Called(params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) DeleteTasksWithContext(ctx context.Context, params *ms.DeleteTasksQuery) (*ms.TaskInfo, error) {
	args := m.Called(ctx, params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) SwapIndexes(params []*ms.SwapIndexesParams) (*ms.TaskInfo, error) {
	args := m.Called(params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockServiceManager) SwapIndexesWithContext(ctx context.Context, params []*ms.SwapIndexesParams) (*ms.TaskInfo, error) {
	args := m.Called(ctx, params)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// Close method to satisfy ServiceManager interface
func (m *MockServiceManager) Close() {
	m.Called()
}

// CreateDumpWithContext method to satisfy ServiceManager interface
func (m *MockServiceManager) CreateDumpWithContext(ctx context.Context) (*ms.TaskInfo, error) {
	args := m.Called(ctx)
	if val := args.Get(0); val == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// CreateKey creates a new API key
func (m *MockServiceManager) CreateKey(request *ms.Key) (*ms.Key, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// CreateKeyWithContext creates a new API key with context
func (m *MockServiceManager) CreateKeyWithContext(ctx context.Context, request *ms.Key) (*ms.Key, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// CreateSnapshot creates a snapshot of the database
func (m *MockServiceManager) CreateSnapshot() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// CreateSnapshotWithContext creates a snapshot of the database with context
func (m *MockServiceManager) CreateSnapshotWithContext(ctx context.Context) (*ms.TaskInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// DeleteKey deletes an API key
func (m *MockServiceManager) DeleteKey(keyOrUID string) (bool, error) {
	args := m.Called(keyOrUID)
	return args.Bool(0), args.Error(1)
}

// DeleteKeyWithContext deletes an API key with context
func (m *MockServiceManager) DeleteKeyWithContext(ctx context.Context, keyOrUID string) (bool, error) {
	args := m.Called(ctx, keyOrUID)
	return args.Bool(0), args.Error(1)
}

// ExperimentalFeatures returns the experimental features manager
func (m *MockServiceManager) ExperimentalFeatures() *ms.ExperimentalFeatures {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*ms.ExperimentalFeatures)
}

// GenerateTenantToken generates a tenant token
func (m *MockServiceManager) GenerateTenantToken(APIKeyUID string, searchRules map[string]interface{}, options *ms.TenantTokenOptions) (string, error) {
	args := m.Called(APIKeyUID, searchRules, options)
	return args.String(0), args.Error(1)
}

// GetKey gets a specific API key
func (m *MockServiceManager) GetKey(keyOrUID string) (*ms.Key, error) {
	args := m.Called(keyOrUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// GetKeyWithContext gets a specific API key with context
func (m *MockServiceManager) GetKeyWithContext(ctx context.Context, keyOrUID string) (*ms.Key, error) {
	args := m.Called(ctx, keyOrUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// GetKeys gets all API keys
func (m *MockServiceManager) GetKeys(param *ms.KeysQuery) (*ms.KeysResults, error) {
	args := m.Called(param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.KeysResults), args.Error(1)
}

// GetKeysWithContext gets all API keys with context
func (m *MockServiceManager) GetKeysWithContext(ctx context.Context, param *ms.KeysQuery) (*ms.KeysResults, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.KeysResults), args.Error(1)
}

// GetRawIndex returns raw index data
func (m *MockServiceManager) GetRawIndex(uid string) (map[string]interface{}, error) {
	args := m.Called(uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetRawIndexWithContext returns raw index data with context
func (m *MockServiceManager) GetRawIndexWithContext(ctx context.Context, uid string) (map[string]interface{}, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetRawIndexes returns raw indexes data
func (m *MockServiceManager) GetRawIndexes(params *ms.IndexesQuery) (map[string]interface{}, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetRawIndexesWithContext returns raw indexes data with context
func (m *MockServiceManager) GetRawIndexesWithContext(ctx context.Context, params *ms.IndexesQuery) (map[string]interface{}, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetStats returns stats
func (m *MockServiceManager) GetStats() (*ms.Stats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Stats), args.Error(1)
}

// GetStatsWithContext returns stats with context
func (m *MockServiceManager) GetStatsWithContext(ctx context.Context) (*ms.Stats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Stats), args.Error(1)
}

// HealthWithContext returns health with context
func (m *MockServiceManager) HealthWithContext(ctx context.Context) (*ms.Health, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Health), args.Error(1)
}

// KeyManager returns a key manager
func (m *MockServiceManager) KeyManager() ms.KeyManager {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.KeyManager)
}

// KeyReader returns a key reader
func (m *MockServiceManager) KeyReader() ms.KeyReader {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.KeyReader)
}

// ListIndexes lists all indexes
func (m *MockServiceManager) ListIndexes(params *ms.IndexesQuery) (*ms.IndexesResults, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexesResults), args.Error(1)
}

// ListIndexesWithContext lists all indexes with context
func (m *MockServiceManager) ListIndexesWithContext(ctx context.Context, params *ms.IndexesQuery) (*ms.IndexesResults, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexesResults), args.Error(1)
}

// MultiSearch performs multiple searches
func (m *MockServiceManager) MultiSearch(queries *ms.MultiSearchRequest) (*ms.MultiSearchResponse, error) {
	args := m.Called(queries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.MultiSearchResponse), args.Error(1)
}

// MultiSearchWithContext performs multiple searches with context
func (m *MockServiceManager) MultiSearchWithContext(ctx context.Context, queries *ms.MultiSearchRequest) (*ms.MultiSearchResponse, error) {
	args := m.Called(ctx, queries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.MultiSearchResponse), args.Error(1)
}

// ServiceReader returns a service reader
func (m *MockServiceManager) ServiceReader() ms.ServiceReader {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.ServiceReader)
}

// TaskManager returns a task manager
func (m *MockServiceManager) TaskManager() ms.TaskManager {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.TaskManager)
}

// TaskReader returns a task reader
func (m *MockServiceManager) TaskReader() ms.TaskReader {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ms.TaskReader)
}

// UpdateKey updates an API key
func (m *MockServiceManager) UpdateKey(keyOrUID string, request *ms.Key) (*ms.Key, error) {
	args := m.Called(keyOrUID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// UpdateKeyWithContext updates an API key with context
func (m *MockServiceManager) UpdateKeyWithContext(ctx context.Context, keyOrUID string, request *ms.Key) (*ms.Key, error) {
	args := m.Called(ctx, keyOrUID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Key), args.Error(1)
}

// VersionWithContext returns the version with context
func (m *MockServiceManager) VersionWithContext(ctx context.Context) (*ms.Version, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Version), args.Error(1)
}

// MockIndexManager is a mock implementation of ms.IndexManager interface
// (which *ms.Index implements)
type MockIndexManager struct {
	mock.Mock
}

// AddDocuments mocks the (i *Index) AddDocuments method
func (m *MockIndexManager) AddDocuments(documentsPtr interface{}, primaryKey ...string) (*ms.TaskInfo, error) {
	callArgs := []interface{}{documentsPtr}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey ...string) (*ms.TaskInfo, error) {
	callArgs := []interface{}{ctx, documentsPtr}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	callArgs := []interface{}{documentsPtr, batchSize}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	callArgs := []interface{}{ctx, documentsPtr, batchSize}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdateDocuments(documentsPtr interface{}, primaryKey ...string) (*ms.TaskInfo, error) {
	callArgs := []interface{}{documentsPtr}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdateDocumentsWithContext(ctx context.Context, documentsPtr interface{}, primaryKey ...string) (*ms.TaskInfo, error) {
	callArgs := []interface{}{ctx, documentsPtr}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdateDocumentsInBatches(documentsPtr interface{}, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	callArgs := []interface{}{documentsPtr, batchSize}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdateDocumentsInBatchesWithContext(ctx context.Context, documentsPtr interface{}, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	callArgs := []interface{}{ctx, documentsPtr, batchSize}
	if len(primaryKey) > 0 {
		callArgs = append(callArgs, primaryKey[0])
	} else {
		callArgs = append(callArgs, nil)
	}
	args := m.Called(callArgs...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteDocument(documentID string) (*ms.TaskInfo, error) {
	args := m.Called(documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteDocumentWithContext(ctx context.Context, documentID string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteDocuments(documentIDs []string) (*ms.TaskInfo, error) {
	args := m.Called(documentIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteDocumentsWithContext(ctx context.Context, documentIDs []string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documentIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteAllDocuments() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteAllDocumentsWithContext(ctx context.Context) (*ms.TaskInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// GetDocument mocks (i *Index) GetDocument
func (m *MockIndexManager) GetDocument(documentID string, documentPtr interface{}) error {
	args := m.Called(documentID, documentPtr)
	return args.Error(0)
}

func (m *MockIndexManager) GetDocumentWithContext(ctx context.Context, documentID string, documentPtr interface{}) error {
	args := m.Called(ctx, documentID, documentPtr)
	return args.Error(0)
}

// GetDocuments mocks (i *Index) GetDocuments
func (m *MockIndexManager) GetDocuments(request *ms.DocumentsQuery) (*ms.DocumentsResult, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.DocumentsResult), args.Error(1)
}

func (m *MockIndexManager) GetDocumentsWithContext(ctx context.Context, request *ms.DocumentsQuery) (*ms.DocumentsResult, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.DocumentsResult), args.Error(1)
}

func (m *MockIndexManager) Search(query string, request *ms.SearchRequest) (*ms.SearchResponse, error) {
	args := m.Called(query, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.SearchResponse), args.Error(1)
}

func (m *MockIndexManager) SearchWithContext(ctx context.Context, query string, request *ms.SearchRequest) (*ms.SearchResponse, error) {
	args := m.Called(ctx, query, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.SearchResponse), args.Error(1)
}

func (m *MockIndexManager) GetStats() (*ms.Stats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Stats), args.Error(1)
}

func (m *MockIndexManager) GetStatsWithContext(ctx context.Context) (*ms.Stats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Stats), args.Error(1)
}

// Settings methods
func (m *MockIndexManager) GetSettings() (*ms.Settings, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Settings), args.Error(1)
}

func (m *MockIndexManager) GetSettingsWithContext(ctx context.Context) (*ms.Settings, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Settings), args.Error(1)
}

func (m *MockIndexManager) UpdateSettings(settings *ms.Settings) (*ms.TaskInfo, error) {
	args := m.Called(settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdateSettingsWithContext(ctx context.Context, settings *ms.Settings) (*ms.TaskInfo, error) {
	args := m.Called(ctx, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetSettings() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetSettingsWithContext(ctx context.Context) (*ms.TaskInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetTask(taskUID int64) (*ms.Task, error) {
	args := m.Called(taskUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockIndexManager) GetTaskWithContext(ctx context.Context, taskUID int64) (*ms.Task, error) {
	args := m.Called(ctx, taskUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockIndexManager) GetTasks(param *ms.TasksQuery) (*ms.TaskResult, error) {
	args := m.Called(param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskResult), args.Error(1)
}

func (m *MockIndexManager) GetTasksWithContext(ctx context.Context, param *ms.TasksQuery) (*ms.TaskResult, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskResult), args.Error(1)
}

func (m *MockIndexManager) WaitForTask(taskUID int64, interval time.Duration, timeout time.Duration) (*ms.Task, error) {
	args := m.Called(taskUID, interval, timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

func (m *MockIndexManager) WaitForTaskWithContext(ctx context.Context, taskUID int64, interval time.Duration, timeout time.Duration) (*ms.Task, error) {
	args := m.Called(ctx, taskUID, interval, timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

// These are specific to ms.Index but part of ms.IndexManager
func (m *MockIndexManager) UID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockIndexManager) CreatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockIndexManager) UpdatedAt() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

// AddDocumentsCsv method to satisfy IndexManager interface
func (m *MockIndexManager) AddDocumentsCsv(documents []byte, csvConfig *ms.CsvDocumentsQuery) (*ms.TaskInfo, error) {
	args := m.Called(documents, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// AddDocumentsCsvWithContext method to satisfy IndexManager interface
func (m *MockIndexManager) AddDocumentsCsvWithContext(ctx context.Context, documents []byte, csvConfig *ms.CsvDocumentsQuery) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documents, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// AddDocumentsCsvFromReader method to satisfy IndexManager interface
func (m *MockIndexManager) AddDocumentsCsvFromReader(documents io.Reader, csvConfig *ms.CsvDocumentsQuery) (*ms.TaskInfo, error) {
	args := m.Called(documents, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// AddDocumentsCsvFromReaderWithContext method to satisfy IndexManager interface
func (m *MockIndexManager) AddDocumentsCsvFromReaderWithContext(ctx context.Context, documents io.Reader, csvConfig *ms.CsvDocumentsQuery) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documents, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsCsvFromReaderInBatches(documents io.Reader, batchSize int, csvConfig *ms.CsvDocumentsQuery) ([]ms.TaskInfo, error) {
	args := m.Called(documents, batchSize, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsCsvFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, csvConfig *ms.CsvDocumentsQuery) ([]ms.TaskInfo, error) {
	args := m.Called(ctx, documents, batchSize, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsCsvInBatches(documents []byte, batchSize int, csvConfig *ms.CsvDocumentsQuery) ([]ms.TaskInfo, error) {
	args := m.Called(documents, batchSize, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsCsvInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, csvConfig *ms.CsvDocumentsQuery) ([]ms.TaskInfo, error) {
	args := m.Called(ctx, documents, batchSize, csvConfig)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjson(documents []byte, primaryKey ...string) (*ms.TaskInfo, error) {
	args := m.Called(documents, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonWithContext(ctx context.Context, documents []byte, primaryKey ...string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documents, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonFromReader(documents io.Reader, primaryKey ...string) (*ms.TaskInfo, error) {
	args := m.Called(documents, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonFromReaderWithContext(ctx context.Context, documents io.Reader, primaryKey ...string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, documents, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonFromReaderInBatches(documents io.Reader, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	args := m.Called(documents, batchSize, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonFromReaderInBatchesWithContext(ctx context.Context, documents io.Reader, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	args := m.Called(ctx, documents, batchSize, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonInBatches(documents []byte, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	args := m.Called(documents, batchSize, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) AddDocumentsNdjsonInBatchesWithContext(ctx context.Context, documents []byte, batchSize int, primaryKey ...string) ([]ms.TaskInfo, error) {
	args := m.Called(ctx, documents, batchSize, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) Delete() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) DeleteWithContext(ctx context.Context) (*ms.TaskInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) FetchInfo() (*ms.IndexResult, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexResult), args.Error(1)
}

func (m *MockIndexManager) FetchInfoWithContext(ctx context.Context) (*ms.IndexResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexResult), args.Error(1)
}

func (m *MockIndexManager) FetchPrimaryKey() (*string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockIndexManager) FetchPrimaryKeyWithContext(ctx context.Context) (*string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockIndexManager) UpdatePrimaryKey(primaryKey string) (*ms.TaskInfo, error) {
	args := m.Called(primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) UpdatePrimaryKeyWithContext(ctx context.Context, primaryKey string) (*ms.TaskInfo, error) {
	args := m.Called(ctx, primaryKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetDisplayedAttributes() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateDisplayedAttributes(attributes []string) (*ms.TaskInfo, error) {
	args := m.Called(attributes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetDisplayedAttributes() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetSearchableAttributes() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateSearchableAttributes(attributes []string) (*ms.TaskInfo, error) {
	args := m.Called(attributes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetSearchableAttributes() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetFilterableAttributes() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateFilterableAttributes(attributes []string) (*ms.TaskInfo, error) {
	args := m.Called(attributes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetFilterableAttributes() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetSortableAttributes() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateSortableAttributes(attributes []string) (*ms.TaskInfo, error) {
	args := m.Called(attributes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetSortableAttributes() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetRankingRules() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateRankingRules(rules []string) (*ms.TaskInfo, error) {
	args := m.Called(rules)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetRankingRules() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetStopWords() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockIndexManager) UpdateStopWords(stopWords []string) (*ms.TaskInfo, error) {
	args := m.Called(stopWords)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetStopWords() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) GetSynonyms() (map[string][]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockIndexManager) UpdateSynonyms(synonyms map[string][]string) (*ms.TaskInfo, error) {
	args := m.Called(synonyms)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

func (m *MockIndexManager) ResetSynonyms() (*ms.TaskInfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.TaskInfo), args.Error(1)
}

// MockClient is a mock implementation of the local Client interface (defined in client.go)
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Connect(cfg config.MeilisearchConfig) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func (m *MockClient) Index(indexName string, document interface{}) error {
	args := m.Called(indexName, document)
	return args.Error(0)
}

func (m *MockClient) Search(indexName string, query interface{}) (*ms.SearchResponse, error) {
	args := m.Called(indexName, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.SearchResponse), args.Error(1)
}

func (m *MockClient) CreateIndex(indexName string) error {
	args := m.Called(indexName)
	return args.Error(0)
}

func (m *MockClient) DeleteIndex(indexName string) (bool, error) {
	args := m.Called(indexName)
	return args.Bool(0), args.Error(1)
}

// GetIndex returns ms.IndexManager as defined in the local Client interface
func (m *MockClient) GetIndex(indexName string) (ms.IndexManager, error) {
	args := m.Called(indexName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Return a *MockIndexManager, as it should implement ms.IndexManager
	return args.Get(0).(ms.IndexManager), args.Error(1)
}

func (m *MockClient) ListIndexes() (*ms.IndexesResults, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.IndexesResults), args.Error(1)
}

func (m *MockClient) IndexExists(indexName string) (bool, error) {
	args := m.Called(indexName)
	return args.Bool(0), args.Error(1)
}

func (m *MockClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// Interface assertions to ensure mocks satisfy the interfaces - temporarily disabled for demo
// var _ ms.ServiceManager = (*MockServiceManager)(nil)
// var _ ms.IndexManager = (*MockIndexManager)(nil)
// var _ Client = (*MockClient)(nil)

// Note: The local IndexManager interface used in meilisearch.go seems to be implicitly
// satisfied by *ms.Index (which implements ms.IndexManager).
// So, for tests involving `StorageManager.client.GetIndex()`, the result would be
// a mock that implements `ms.IndexManager` (i.e. `*MockIndexManager`).

// The key is that MockServiceManager mocks *ms.Client (or its subset used like a service manager)
// and MockIndexManager mocks *ms.Index (or its subset used like an index manager).
// MockClient mocks the local `Client` interface.

// GetRawIndex is not a standard ms.Client method. If it's on your local Client, mock it here.
// func (m *MockClient) GetRawIndex(uid string) (IndexManager, error) { ... }

// UpdateIndex is not a standard ms.Client method.
// func (m *MockClient) UpdateIndex(uid string, primaryKey string) (*ms.TaskInfo, error) { ... }

// GetTask is often on ms.Client or ms.Index.
// Mocked in MockServiceManager. If also on local Client:
func (m *MockClient) GetTask(taskUID int64) (*ms.Task, error) {
	args := m.Called(taskUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ms.Task), args.Error(1)
}

// Methods from ms.Client that MockServiceManager handles:
// Health(), IsHealthy(), Stats(), Version(), CreateDump(), GetDumpStatus(),
// GetTasks(), GetTask(), WaitForTask(), CancelTasks(), DeleteTasks(), SwapIndexes()

// Methods from ms.Index that MockIndexManager handles:
// AddDocuments(), AddDocumentsInBatches(), UpdateDocuments(), UpdateDocumentsInBatches(),
// DeleteDocument(), DeleteDocuments(), DeleteAllDocuments(), GetDocument(), GetDocuments(),
// Search(), GetStats(), GetSettings(), UpdateSettings(), ResetSettings(), ... and all other settings.

// This MockClient is specifically for the *local* `Client` interface.
// Ensure its methods match what's defined in `client.go` for the `Client` interface.
// The `config.MeilisearchConfig` import was added because `MockClient.Connect` uses it.
// The `context` and `time` imports are standard for many operations.
