package campaign

import (
	"trinity/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateCampaign(campaign *model.Campaign) (string, error) {
	args := m.Called(campaign)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetCampaignByID(id string) (*model.Campaign, error) {
	args := m.Called(id)
	if campaign, ok := args.Get(0).(*model.Campaign); ok {
		return campaign, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) IncrementUsedUsers(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) ListCampaigns() ([]model.Campaign, error) {
	args := m.Called()
	if campaigns, ok := args.Get(0).([]model.Campaign); ok {
		return campaigns, args.Error(1)
	}
	return nil, args.Error(1)
}
