package voucher

import (
	"trinity/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) RedeemVoucher(code string, userID string) (*model.Voucher, error) {
	args := m.Called(code, userID)
	voucher := args.Get(0)
	if voucher == nil {
		return nil, args.Error(1)
	}
	return voucher.(*model.Voucher), args.Error(1)
}
