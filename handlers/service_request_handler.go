package handlers

import (
	. "clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	serviceRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "create_service_request_handler_counter",
		Help: "The total number of service requests created",
	}, []string{"workflow_name"})
	serviceRequestHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "create_service_request_handler_histogram",
		Help: "The total number of service requests created",
	})
)
// Create Service Request godoc
// @Summary Create a service request
// @Description Create a service request and get service request id
// @Accept json
// @Produce json
// @Param workflowname path string true "Workflow Name"
// @Param serviceRequestPayload body string true "Service Request Payload"
// @Success 200 {object} models.ServiceRequestResponse
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /serviceRequest/{workflowname} [post]
func createServiceRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()
		log.Println("Create service request handler")
		workflowName := c.Param("workflowName")
		serviceRequestCounter.WithLabelValues(workflowName).Inc()
		_, err := services.FindWorkflowByName(workflowName)

		requestPayload := readRequestPayload(c)

		if err != nil {
			errorResponse := CreateErrorResponse(http.StatusBadRequest, "No record found with given workflow name : "+workflowName)
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		// Create new service request
		serviceReq := NewServiceRequest(workflowName, requestPayload)
		serviceReq, err = services.SaveServiceRequest(serviceReq)
		if err != nil {
			errorResponse := CreateErrorResponse(http.StatusBadRequest, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		services.AddServiceRequestToChannel(serviceReq)
		response := prepareServiceRequestResponse(serviceReq)
		serviceRequestHistogram.Observe(time.Since(startTime).Seconds())
		c.JSON(http.StatusOK, response)
	}
}

func prepareServiceRequestResponse(serviceReq ServiceRequest) ServiceRequestResponse {
	response := ServiceRequestResponse{
		URL:    "/serviceRequest/" + serviceReq.ID.String(),
		Status: serviceReq.Status,
		ID:     serviceReq.ID,
	}
	return response
}

func readRequestPayload(c *gin.Context) map[string]interface{} {
	var payload map[string]interface{}
	if c.Request.Body != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(data, &payload)
		log.Println("Request Body", payload)
		return payload
	} else {
		return nil
	}
}
// Get Service Request By Id godoc
// @Summary Get service request details by service request id
// @Description Get service request by service request id
// @Accept json
// @Produce json
// @Param serviceRequestId path string true "Service Request Id"
// @Success 200 {object} models.ServiceRequestStatusResponse
// @Failure 400 {object} models.ClampErrorResponse
// @Failure 404 {object} models.ClampErrorResponse
// @Failure 500 {object} models.ClampErrorResponse
// @Router /serviceRequest/{serviceRequestId} [get]
func getServiceRequestStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceRequestId := c.Param("serviceRequestId")

		serviceRequest, err := services.FindServiceRequestByID(uuid.MustParse(serviceRequestId))
		if err != nil {
			c.JSON(http.StatusBadRequest, CreateErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		workflow, _ := services.FindWorkflowByName(serviceRequest.WorkflowName)
		stepsStatues, _ := services.FindStepStatusByServiceRequestId(uuid.MustParse(serviceRequestId))
		stepsStatusResponse := services.PrepareStepStatusResponse(uuid.MustParse(serviceRequestId), workflow, stepsStatues)
		//TODO - handle error scenario. Currently it is always 200 ok
		c.JSON(http.StatusOK, stepsStatusResponse)
	}
}
