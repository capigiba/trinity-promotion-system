package campaign

import (
	"errors"
	"testing"
	"time"
	"trinity/internal/model"
	"trinity/internal/voucher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Initialize the service with mocked dependencies
func setupService(mockRepo *MockRepository, mockVoucherRepo *voucher.MockRepository) *service {
	return &service{
		repo:        mockRepo,
		voucherRepo: mockVoucherRepo,
	}
}

func TestService_CreateCampaign_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	campaign := &model.Campaign{
		Name:        "Test Campaign",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(48 * time.Hour),
		MaxUsers:    100,
		UsedUsers:   0,
		Description: "A test campaign",
	}

	mockRepo.On("CreateCampaign", campaign).Return("campaign123", nil)

	id, err := service.CreateCampaign(campaign)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "campaign123", id, "Expected campaign ID to match")
	mockRepo.AssertExpectations(t)
}

func TestService_CreateCampaign_InvalidDates(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	campaign := &model.Campaign{
		Name:        "Invalid Campaign",
		StartDate:   time.Now().Add(48 * time.Hour),
		EndDate:     time.Now(),
		MaxUsers:    100,
		UsedUsers:   0,
		Description: "Campaign with invalid dates",
	}

	id, err := service.CreateCampaign(campaign)

	assert.Error(t, err, "Expected an error due to invalid dates")
	assert.Equal(t, "", id, "Expected no campaign ID to be returned")

	mockRepo.AssertNotCalled(t, "CreateCampaign", mock.Anything)
}

func TestService_GenerateVouchers_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	campaignID := "campaign123"
	count := 5

	campaign := &model.Campaign{
		Id:          campaignID,
		Name:        "Voucher Campaign",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(72 * time.Hour),
		MaxUsers:    10,
		UsedUsers:   3,
		Description: "Campaign for generating vouchers",
	}

	// Expect GetCampaignByID to be called and return the campaign
	mockRepo.On("GetCampaignByID", campaignID).Return(campaign, nil)

	// Expect IncrementUsedUsers to be called once
	mockRepo.On("IncrementUsedUsers", campaignID).Return(nil)

	// Expect CreateVoucher to be called 'count' times with any Voucher
	mockVoucherRepo.On("CreateVoucher", mock.AnythingOfType("*model.Voucher")).Return(nil).Times(count)

	vouchers, err := service.GenerateVouchers(campaignID, count)

	assert.NoError(t, err, "Expected no error")
	assert.Len(t, vouchers, count, "Expected number of generated vouchers to match")

	mockRepo.AssertExpectations(t)
	mockVoucherRepo.AssertExpectations(t)
}

func TestService_GenerateVouchers_NotEnoughVouchers(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	campaignID := "campaign123"
	count := 8

	campaign := &model.Campaign{
		Id:          campaignID,
		Name:        "Limited Voucher Campaign",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(72 * time.Hour),
		MaxUsers:    10,
		UsedUsers:   5, // Only 5 vouchers left
		Description: "Campaign with limited vouchers",
	}

	// Expect GetCampaignByID to be called and return the campaign
	mockRepo.On("GetCampaignByID", campaignID).Return(campaign, nil)

	vouchers, err := service.GenerateVouchers(campaignID, count)

	assert.Error(t, err, "Expected an error due to insufficient vouchers")
	assert.Nil(t, vouchers, "Expected no vouchers to be returned")

	mockRepo.AssertExpectations(t)
	mockVoucherRepo.AssertNotCalled(t, "CreateVoucher", mock.Anything)
	mockRepo.AssertNotCalled(t, "IncrementUsedUsers", campaignID)
}

func TestService_ListCampaigns_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	campaigns := []model.Campaign{
		{
			Id:          "campaign1",
			Name:        "Campaign One",
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(48 * time.Hour),
			MaxUsers:    100,
			UsedUsers:   20,
			Description: "First campaign",
		},
		{
			Id:          "campaign2",
			Name:        "Campaign Two",
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(72 * time.Hour),
			MaxUsers:    200,
			UsedUsers:   50,
			Description: "Second campaign",
		},
	}

	mockRepo.On("ListCampaigns").Return(campaigns, nil)

	result, err := service.ListCampaigns()

	assert.NoError(t, err, "Expected no error")
	assert.Len(t, result, 2, "Expected two campaigns")
	assert.Equal(t, campaigns, result, "Campaign lists should match")
	mockRepo.AssertExpectations(t)
}

func TestService_ListCampaigns_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockVoucherRepo := new(voucher.MockRepository)
	service := setupService(mockRepo, mockVoucherRepo)

	mockRepo.On("ListCampaigns").Return(nil, errors.New("database error"))

	result, err := service.ListCampaigns()

	assert.Error(t, err, "Expected an error from repository")
	assert.Nil(t, result, "Expected no campaigns to be returned")
	mockRepo.AssertExpectations(t)
}
