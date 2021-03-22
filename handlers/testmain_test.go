package handlers

import (
	"clamp-core/config"
	"clamp-core/executors"
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/services"
	"clamp-core/transform"
	"clamp-core/utils"
	"fmt"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

const testWorkflowName string = "testWorkflow"
const testTransformationWorkflow string = "testTransformationWorkflow"

var mockDB repository.MockDB
var testHTTRouter *gin.Engine
var testHTTPServer *httptest.Server

func TestMain(m *testing.M) {
	err := config.Load()
	if err != nil {
		fmt.Printf("Loading config failed: %s\n", err)
	}

	repository.SetDB(repository.NewMemoryDB())

	err = services.InitServiceRequestWorkers()
	if err != nil {
		fmt.Printf("Initializinng service request workers failed: %s", err)
	}

	err = services.InitResumeWorkers()
	if err != nil {
		fmt.Printf("Initializinng resume workers failed: %s", err)
	}

	gin.SetMode(gin.TestMode)

	testHTTRouter = setupRouter()

	testHTTPServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpResponseBody := map[string]interface{}{"id": "1234", "name": "ABC", "email": "abc@sahaj.com", "org": "sahaj"}
		json.NewEncoder(w).Encode(httpResponseBody)
	}))

	step := models.Step{
		Name:      "1",
		Type:      utils.StepTypeSync,
		Mode:      utils.StepModeHTTP,
		Transform: false,
		Enabled:   false,
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     testHTTPServer.URL,
			Headers: "",
		},
	}

	workflow := models.Workflow{
		Name:  testWorkflowName,
		Steps: []models.Step{step},
	}

	_, err = services.SaveWorkflow(&workflow)
	if err != nil {
		panic(err)
	}

	step = models.Step{
		Name:      "1",
		Type:      utils.StepTypeSync,
		Mode:      utils.StepModeHTTP,
		Transform: true,
		Enabled:   false,
		RequestTransform: &transform.JSONTransform{
			Spec: map[string]interface{}{"name": "test"},
		},
		Val: &executors.HTTPVal{
			Method:  "POST",
			URL:     testHTTPServer.URL,
			Headers: "",
		},
	}

	workflow = models.Workflow{
		Name:  testTransformationWorkflow,
		Steps: []models.Step{step},
	}

	_, err = services.SaveWorkflow(&workflow)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
