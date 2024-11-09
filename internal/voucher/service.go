package voucher

import (
	"errors"
	"time"
	"trinity/internal/model"
	"trinity/pkg/logger"
)

// Service defines voucher business logic methods
type Service interface {
	RedeemVoucher(code string, userID string) (*model.Voucher, error)
}

// service implements Service interface
type service struct {
	repo   Repository
	logger logger.Logger
}

// NewService creates a new Voucher service
func NewService(repo Repository) Service {
	return &service{
		repo:   repo,
		logger: logger.NewLogger("voucherService"),
	}
}

// RedeemVoucher redeems a voucher
func (s *service) RedeemVoucher(code string, userID string) (*model.Voucher, error) {
	voucher, err := s.repo.GetVoucherByCode(code)
	if err != nil {
		s.logger.Errorf("Failed to get voucher: %v", err)
		return nil, errors.New("invalid voucher code")
	}

	if voucher.Used {
		return nil, errors.New("voucher already used")
	}

	if time.Now().After(voucher.ExpiryDate) {
		return nil, errors.New("voucher expired")
	}

	// Mark voucher as used
	voucher.Used = true
	voucher.UserId = userID

	err = s.repo.UpdateVoucher(voucher)
	if err != nil {
		s.logger.Errorf("Failed to update voucher: %v", err)
		return nil, errors.New("failed to update voucher")
	}

	return voucher, nil
}
