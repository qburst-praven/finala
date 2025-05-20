package api

import (
	"encoding/json"
	"finala/api/config"
	"finala/api/email_utility"
	"finala/api/httpparameters"
	"finala/api/storage"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	queryParamFilterPrefix     = "filter_"
	resourceTrendsLimitDefault = 60
)

// DetectEventsInfo describes the incoming HTTP events
type DetectEventsInfo struct {
	ResourceName string
	EventType    string
	EventTime    int64
	Data         interface{}
}

type ReportAPIResponse struct {
	Message string     `json:"message"`
	Status  int        `json:"status"`
	Data    [][]string `json:"data"`
}

// GetSummary return list of summary executions
func (server *Server) GetSummary(resp http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	executionID := req.PathValue("executionID")
	filters := httpparameters.GetFilterQueryParamWithOutPrefix(queryParamFilterPrefix, queryParams)

	response, err := server.storage.GetSummary(executionID, filters)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return

	}
	server.JSONWrite(resp, http.StatusOK, response)
}

// GetExecutions return list collector executions
func (server *Server) GetExecutions(resp http.ResponseWriter, req *http.Request) {
	querylimit, _ := strconv.Atoi(httpparameters.QueryParamWithDefault(req, "querylimit", storage.GetExecutionsQueryLimit))
	results, err := server.storage.GetExecutions(querylimit)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return

	}
	server.JSONWrite(resp, http.StatusOK, results)
}

// GetResourceData return resuts details by resource type
func (server *Server) GetResourceData(resp http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	queryErrs := url.Values{}
	resourceType := req.PathValue("type")
	var search string
	filters := httpparameters.GetFilterQueryParamWithOutPrefix(queryParamFilterPrefix, queryParams)

	executionID := req.URL.Query().Get("executionID")
	if executionID == "" {
		queryErrs.Add("executionID", "executionID field is mandatory")
	}

	if len(queryErrs) > 0 {
		server.JSONWrite(resp, http.StatusBadRequest, HttpErrorResponse{ErrorQuery: queryErrs})
		return
	}

	response, err := server.storage.GetResources(resourceType, executionID, filters, search)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return

	}
	server.JSONWrite(resp, http.StatusOK, response)
}

// GetResourceTrends return trends by resource type, id, region and metric
func (server *Server) GetResourceTrends(resp http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	resourceType := req.PathValue("type")
	filters := httpparameters.GetFilterQueryParamWithOutPrefix(queryParamFilterPrefix, queryParams)

	limitString := req.URL.Query().Get("limit")
	var limit int = resourceTrendsLimitDefault
	var err error
	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil || limit < 1 {
			limit = resourceTrendsLimitDefault
		}
	}

	trends, err := server.storage.GetResourceTrends(resourceType, filters, limit)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return

	}
	server.JSONWrite(resp, http.StatusOK, trends)
}

// GetExecutionTags return resuts details by resource type
func (server *Server) GetExecutionTags(resp http.ResponseWriter, req *http.Request) {
	executionID := req.PathValue("executionID")

	response, err := server.storage.GetExecutionTags(executionID)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return

	}
	server.JSONWrite(resp, http.StatusOK, response)
}

// DetectEvents save collectors events data
func (server *Server) DetectEvents(resp http.ResponseWriter, req *http.Request) {
	executionID := req.PathValue("executionID")

	buf, bodyErr := io.ReadAll(req.Body)

	if bodyErr != nil {
		server.JSONWrite(resp, http.StatusBadRequest, HttpErrorResponse{Error: bodyErr.Error()})
		return
	}

	var detectEventsInfo []DetectEventsInfo
	err := json.Unmarshal(buf, &detectEventsInfo)
	if err != nil {
		server.JSONWrite(resp, http.StatusBadRequest, HttpErrorResponse{Error: err.Error()})
		return
	}

	log.WithFields(log.Fields{
		"events": len(detectEventsInfo),
	}).Info("Got bulk events")

	go func() {
		for _, event := range detectEventsInfo {

			rowData := storage.EventRow{
				ExecutionID:  executionID,
				ResourceName: event.ResourceName,
				EventType:    event.EventType,
				EventTime:    event.EventTime,
				Timestamp:    time.Now(),
				Data:         event.Data,
			}
			bolB, _ := json.Marshal(rowData)
			server.storage.Save(string(bolB))
		}
	}()

	server.JSONWrite(resp, http.StatusAccepted, nil)
}

// NotFoundRoute return when route not found
func (server *Server) NotFoundRoute(resp http.ResponseWriter, req *http.Request) {
	server.JSONWrite(resp, http.StatusNotFound, HttpErrorResponse{Error: "Path not found"})
}

// HealthCheckHandler return ok if server is up
func (server *Server) HealthCheckHandler(resp http.ResponseWriter, req *http.Request) {
	server.JSONWrite(resp, http.StatusOK, HealthResponse{Status: true})
}

// VersionHandler returns the latest Finala version
func (server *Server) VersionHandler(resp http.ResponseWriter, req *http.Request) {
	version, err := server.version.Get()
	if err != nil {
		server.JSONWrite(resp, http.StatusNotFound, HttpErrorResponse{Error: "Version was not found"})
		return
	}
	server.JSONWrite(resp, http.StatusOK, version)
}

// send pdf report via mail
func (server *Server) SendReport(resp http.ResponseWriter, req *http.Request) {

	buf, bodyErr := io.ReadAll(req.Body)

	if bodyErr != nil {
		server.JSONWrite(resp, http.StatusBadRequest, HttpErrorResponse{Error: bodyErr.Error()})
		return
	}

	var sendEmailInfo config.SendEmailInfo
	err := json.Unmarshal(buf, &sendEmailInfo)
	if err != nil {
		server.JSONWrite(resp, http.StatusBadRequest, HttpErrorResponse{Error: err.Error()})
		return
	}

	toEmails := sendEmailInfo.ToEmails
	executionID := sendEmailInfo.ExecutionID
	resourceType := sendEmailInfo.ResourceType

	if executionID == "" {
		server.JSONWrite(resp, http.StatusOK, "Execution Id is mandatory")
		return
	}
	if resourceType == "" {
		server.JSONWrite(resp, http.StatusOK, "Resource Type is mandatory")
		return
	}
	responseMsg := "Email sent successfully"
	statusCode := 200

	//queryParams := req.URL.Query()
	// filters := map[string]string{}
	//filters := httpparameters.GetFilterQueryParamWithOutPrefix(queryParamFilterPrefix, queryParams)
	responseData, err := server.storage.GetResources(resourceType, executionID, sendEmailInfo.Filters, sendEmailInfo.Search)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: err.Error()})
		return
	}

	if len(responseData) == 0 {
		server.JSONWrite(resp, http.StatusOK, ReportAPIResponse{Message: "No data", Status: statusCode})
		return
	}

	pdfFileName := "Finala_report.pdf"
	pdfContent := "A Comprehensive Analysis of Efficiency Factors and Recommendations for Improvement"
	email_utility.CreatePDF(pdfFileName, pdfContent, responseData, sendEmailInfo)
	emailConfig, err := config.LoadAPI("/etc/finala/config.yaml")
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: "Failed to load email configuration: " + err.Error()})
		return
	}
	username := emailConfig.SMTPConf.Username
	password := emailConfig.SMTPConf.Password
	smtpServer := emailConfig.SMTPConf.SMTPServer
	smtpPort, err := strconv.Atoi(emailConfig.SMTPConf.SMTPPort)
	if err != nil {
		server.JSONWrite(resp, http.StatusInternalServerError, HttpErrorResponse{Error: "Invalid SMTP port in configuration: " + err.Error()})
		return
	}
	subject := "Finala Report"
	body := "<p>Kindly review the attached PDF for the comprehensive report on Finala.</p>"
	sender := email_utility.NewSMTPSender(smtpServer, smtpPort, username, password)
	err = sender.Send(toEmails, subject, body, pdfFileName)
	if err != nil {
		responseMsg = "Error in sending mail"
		statusCode = 500
		log.WithFields(log.Fields{"events": len("test")}).Info("---------", err)
	}

	server.JSONWrite(resp, http.StatusOK, ReportAPIResponse{Message: responseMsg, Status: statusCode})
}
