package voucher

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
)

// setupRouter initializes the Gin engine with the voucher routes
func setupRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	voucherGroup := r.Group("/vouchers")
	handler.RegisterRoutes(voucherGroup)
	return r
}

// TestRedeemVoucher_Success tests successful voucher redemption
func TestRedeemVoucher_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)
	router := setupRouter(handler)

	reqBody := RedeemVoucherRequest{
		Code:   "VALIDCODE",
		UserId: "user123",
	}
	body, _ := json.Marshal(reqBody)

	// Create a voucher with a valid ObjectID
	voucher := &model.Voucher{
		Id:         "TEST",
		Code:       "VALIDCODE",
		CampaignID: "campaign123",
		Used:       true,
		UserId:     "user123",
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}

	// Mock the service's RedeemVoucher method
	mockService.On("RedeemVoucher", "VALIDCODE", "user123").Return(voucher, nil)

	// Create a new HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/vouchers/redeem", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

	var responseVoucher model.Voucher
	err := json.Unmarshal(w.Body.Bytes(), &responseVoucher)
	assert.NoError(t, err, "Response should be a valid Voucher")

	assert.Equal(t, voucher.Id, responseVoucher.Id, "Voucher IDs should match")
	assert.Equal(t, voucher.Code, responseVoucher.Code, "Voucher codes should match")
	assert.Equal(t, voucher.CampaignID, responseVoucher.CampaignID, "CampaignIDs should match")
	assert.Equal(t, voucher.Used, responseVoucher.Used, "Voucher used status should match")
	assert.Equal(t, voucher.UserId, responseVoucher.UserId, "UserIds should match")
	assert.WithinDuration(t, voucher.ExpiryDate, responseVoucher.ExpiryDate, time.Second, "ExpiryDate should be approximately equal")

	// Ensure that the mock was called as expected
	mockService.AssertExpectations(t)
}
