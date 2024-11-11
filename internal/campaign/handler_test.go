package campaign

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trinity/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SetupHandler initializes the handler with mocked dependencies
func SetupHandler(mockService *MockService) *Handler {
	return NewHandler(mockService)
}

// Helper function to perform HTTP requests
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody bytes.Buffer
	if body != nil {
		json.NewEncoder(&reqBody).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHandler_CreateCampaign_Success(t *testing.T) {
	mockService := new(MockService)
	handler := SetupHandler(mockService)

	router := gin.Default()
	router.POST("/campaigns", handler.CreateCampaign)

	requestBody := CreateCampaignRequest{
		Name:        "Test Campaign",
		Discount:    20.5,
		MaxUsers:    100,
		Description: "A test campaign",
		StartDate:   time.Now().Format(time.RFC3339),
		EndDate:     time.Now().Add(48 * time.Hour).Format(time.RFC3339),
	}

	campaignID := "campaign123"

	// Set expectations
	mockService.On("CreateCampaign", mock.AnythingOfType("*model.Campaign")).Return(campaignID, nil)

	w := performRequest(router, "POST", "/campaigns", requestBody)

	assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")
	var response model.Campaign
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "Expected no error unmarshaling response")
	assert.Equal(t, campaignID, response.Id, "Expected campaign ID to match")
	assert.Equal(t, requestBody.Name, response.Name, "Expected campaign name to match")
	assert.Equal(t, requestBody.Discount, response.Discount, "Expected campaign discount to match")
	assert.Equal(t, requestBody.MaxUsers, response.MaxUsers, "Expected campaign max_users to match")
	assert.Equal(t, requestBody.Description, response.Description, "Expected campaign description to match")

	mockService.AssertExpectations(t)
}
