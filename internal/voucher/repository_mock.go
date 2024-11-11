package voucher

import (
	"trinity/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateVoucher(voucher *model.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}

func (m *MockRepository) GetVoucherByCode(code string) (*model.Voucher, error) {
	args := m.Called(code)
	voucher := args.Get(0)
	if voucher == nil {
		return nil, args.Error(1)
	}
	return voucher.(*model.Voucher), args.Error(1)
}

func (m *MockRepository) UpdateVoucher(voucher *model.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}
