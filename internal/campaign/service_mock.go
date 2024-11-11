package campaign

import (
	"trinity/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateCampaign(campaign *model.Campaign) (string, error) {
	args := m.Called(campaign)
	return args.String(0), args.Error(1)
}

func (m *MockService) GenerateVouchers(campaignID string, count int) ([]model.Voucher, error) {
	args := m.Called(campaignID, count)
	return args.Get(0).([]model.Voucher), args.Error(1)
}

func (m *MockService) ListCampaigns() ([]model.Campaign, error) {
	args := m.Called()
	return args.Get(0).([]model.Campaign), args.Error(1)
}
