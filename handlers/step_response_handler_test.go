package handlers

import (
	"bytes"
	"clamp-core/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecordStepResponse(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w := httptest.NewRecorder()
	res := models.AsyncStepResponse{
		ServiceRequestID: uuid.UUID{},
		StepID:           0,
		Response:         nil,
	}
	workflowJSONReg, _ := json.Marshal(res)
	requestReader := bytes.NewReader(workflowJSONReg)

	req, _ := http.NewRequest("POST", "/stepResponse", requestReader)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampSuccessResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("success", jsonResp.Message)
}

func TestShouldReturnBadRequestWhenRequestContainsInvalidDataForRecordStepResponse(t *testing.T) {
	assert := assert.New(t)
	setUpFixture()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/stepResponse", nil)
	testHTTRouter.ServeHTTP(w, req)

	bodyStr := w.Body.String()
	var jsonResp models.ClampErrorResponse
	json.Unmarshal([]byte(bodyStr), &jsonResp)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("invalid request", jsonResp.Message)
}
