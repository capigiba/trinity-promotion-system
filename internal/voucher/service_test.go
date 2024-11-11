package voucher

import (
	"errors"
	"testing"
	"time"
	"trinity/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestServiceRedeemVoucher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	code := "VALIDCODE"
	userID := "user123"

	voucher := &model.Voucher{
		Code:       code,
		Used:       false,
		UserId:     "",
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}

	// Mock GetVoucherByCode
	mockRepo.On("GetVoucherByCode", code).Return(voucher, nil)

	// Mock UpdateVoucher
	updatedVoucher := *voucher
	updatedVoucher.Used = true
	updatedVoucher.UserId = userID
	mockRepo.On("UpdateVoucher", &updatedVoucher).Return(nil)

	result, err := service.RedeemVoucher(code, userID)
	assert.NoError(t, err, "Redeeming a valid voucher should not return an error")
	assert.Equal(t, true, result.Used, "Voucher should be marked as used")
	assert.Equal(t, userID, result.UserId, "UserId should be updated")

	mockRepo.AssertExpectations(t)
}

func TestServiceRedeemVoucher_InvalidCode(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	code := "INVALIDCODE"
	userID := "user123"

	// Mock GetVoucherByCode to return error
	mockRepo.On("GetVoucherByCode", code).Return(nil, errors.New("voucher not found"))

	result, err := service.RedeemVoucher(code, userID)
	assert.Error(t, err, "Redeeming with invalid code should return an error")
	assert.Nil(t, result, "Result should be nil for invalid code")

	mockRepo.AssertExpectations(t)
}

func TestServiceRedeemVoucher_AlreadyUsed(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	code := "USEDVOUCHER"
	userID := "user123"

	voucher := &model.Voucher{
		Code:       code,
		Used:       true,
		UserId:     "anotherUser",
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}

	// Mock GetVoucherByCode
	mockRepo.On("GetVoucherByCode", code).Return(voucher, nil)

	result, err := service.RedeemVoucher(code, userID)
	assert.Error(t, err, "Redeeming an already used voucher should return an error")
	assert.Nil(t, result, "Result should be nil for already used voucher")

	mockRepo.AssertExpectations(t)
}

func TestServiceRedeemVoucher_Expired(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	code := "EXPIREDVOUCHER"
	userID := "user123"

	voucher := &model.Voucher{
		Code:       code,
		Used:       false,
		UserId:     "",
		ExpiryDate: time.Now().Add(-1 * time.Hour), // Expired
	}

	// Mock GetVoucherByCode
	mockRepo.On("GetVoucherByCode", code).Return(voucher, nil)

	result, err := service.RedeemVoucher(code, userID)
	assert.Error(t, err, "Redeeming an expired voucher should return an error")
	assert.Nil(t, result, "Result should be nil for expired voucher")

	mockRepo.AssertExpectations(t)
}

func TestServiceRedeemVoucher_UpdateError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	code := "UPDATEERROR"
	userID := "user123"

	voucher := &model.Voucher{
		Code:       code,
		Used:       false,
		UserId:     "",
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}

	// Mock GetVoucherByCode
	mockRepo.On("GetVoucherByCode", code).Return(voucher, nil)

	// Mock UpdateVoucher to return error
	updatedVoucher := *voucher
	updatedVoucher.Used = true
	updatedVoucher.UserId = userID
	mockRepo.On("UpdateVoucher", &updatedVoucher).Return(errors.New("update failed"))

	result, err := service.RedeemVoucher(code, userID)
	assert.Error(t, err, "Redeeming should return an error if update fails")
	assert.Nil(t, result, "Result should be nil if update fails")

	mockRepo.AssertExpectations(t)
}
