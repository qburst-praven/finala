package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	log "github.com/sirupsen/logrus"

	authhandlers "finala/api/handlers"
	"finala/api/storage"
	"finala/serverutil"
	"finala/version"
)

const (
	// DrainTimeout is how long to wait until the server is drained before closing it
	DrainTimeout = time.Second * 30
)

// Server is the API server struct
type Server struct {
	router     *http.ServeMux
	httpserver *http.Server
	storage    storage.StorageDescriber
	version    version.VersionManagerDescriptor
}

// NewServer returns a new Server
func NewServer(port int, storage storage.StorageDescriber, version version.VersionManagerDescriptor) *Server {

	router := http.NewServeMux()
	// Define more specific CORS options
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"})

	return &Server{
		router:  router,
		storage: storage,
		version: version,
		httpserver: &http.Server{
			// Apply the more specific CORS options
			Handler: handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router),
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		},
	}
}

// Serve starts the HTTP server and listens until StopFunc is called
func (server *Server) Serve() serverutil.StopFunc {
	ctx, cancelFn := context.WithCancel(context.Background())
	server.BindEndpoints()

	stopped := make(chan bool)
	go func() {
		<-ctx.Done()
		serverCtx, serverCancelFn := context.WithTimeout(context.Background(), DrainTimeout)
		err := server.httpserver.Shutdown(serverCtx)
		if err != nil {
			log.WithError(err).Error("error occurred while shutting down manager HTTP server")
		}
		serverCancelFn()
		stopped <- true
	}()
	go func() {
		log.WithField("address", server.httpserver.Addr).Info("server listening on")
		err := server.httpserver.ListenAndServe()
		if err != nil {
			log.WithError(err).Info("HTTP server status")
		}
	}()
	return func() {
		cancelFn()
		<-stopped
		log.Warn("HTTP server has been drained and shut down")
	}
}

// BindEndpoints sets up the router to handle API endpoints
func (server *Server) BindEndpoints() {
	// Add pattern handlers using Go 1.22's ServeMux
	server.router.HandleFunc("GET /api/v1/summary/{executionID}", server.GetSummary)
	server.router.HandleFunc("GET /api/v1/executions", server.GetExecutions)
	server.router.HandleFunc("GET /api/v1/resources/{type}", server.GetResourceData)
	server.router.HandleFunc("GET /api/v1/trends/{type}", server.GetResourceTrends)
	server.router.HandleFunc("GET /api/v1/tags/{executionID}", server.GetExecutionTags)
	server.router.HandleFunc("POST /api/v1/detect-events/{executionID}", server.DetectEvents)
	server.router.HandleFunc("POST /api/v1/send-report", server.SendReport)
	server.router.HandleFunc("GET /api/v1/version", server.VersionHandler)
	server.router.HandleFunc("GET /api/v1/health", server.HealthCheckHandler)

	// ADDED: Login route
	server.router.HandleFunc("POST /api/v1/auth/login", authhandlers.LoginHandler)

	// Add a catch-all handler for not found routes
	server.router.HandleFunc("/", server.NotFoundRoute)
}

// Router returns the Go ServeMux HTTP router defined for this server
func (server *Server) Router() *http.ServeMux {
	return server.router
}

// JSONWrite return JSON response to the client
func (server *Server) JSONWrite(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	encoder := json.NewEncoder(resp)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(data)
	if err != nil {
		log.WithError(err).Error("could not set message error in json response")
	}
}
