package purchase

import (
	"errors"
	"time"
	"trinity/internal/model"
	"trinity/internal/subscription"
	"trinity/internal/voucher"
	"trinity/pkg/logger"
)

// Service defines purchase business logic methods
type Service interface {
	ProcessPurchase(userID string, plan model.SubscriptionPlan, voucherCode string) (*model.Purchase, error)
}

// service implements Service interface
type service struct {
	purchaseRepo     Repository
	voucherRepo      voucher.Repository
	subscriptionRepo subscription.Repository
	logger           logger.Logger
}

// NewService creates a new Purchase service
func NewService(purchaseRepo Repository, voucherRepo voucher.Repository, subscriptionRepo subscription.Repository) Service {
	return &service{
		purchaseRepo:     purchaseRepo,
		voucherRepo:      voucherRepo,
		subscriptionRepo: subscriptionRepo,
		logger:           logger.NewLogger("purchaseService"),
	}
}

// ProcessPurchase processes a purchase
func (s *service) ProcessPurchase(userId string, plan model.SubscriptionPlan, voucherCode string) (*model.Purchase, error) {
	// Define base prices
	var basePrice float64
	switch plan {
	case model.PlanSilver:
		basePrice = 100.0 // Base price for Silver plan
	case model.PlanGold:
		basePrice = 200.0 // Base price for Gold plan
	default:
		return nil, errors.New("invalid subscription plan")
	}

	// Initialize discount
	discount := 0.0

	// If voucher code is provided, validate and apply discount
	if voucherCode != "" {
		voucher, err := s.voucherRepo.GetVoucherByCode(voucherCode)
		if err != nil {
			return nil, errors.New("invalid voucher code")
		}
		if voucher.Used {
			return nil, errors.New("voucher already used")
		}
		if time.Now().After(voucher.ExpiryDate) {
			return nil, errors.New("voucher expired")
		}

		// Apply discount (assuming 30% as per requirement)
		discount = basePrice * 0.30

		// Mark voucher as used
		voucher.Used = true
		voucher.UserId = userId
		err = s.voucherRepo.UpdateVoucher(voucher)
		if err != nil {
			return nil, errors.New("failed to update voucher")
		}
	}

	// Calculate total amount
	totalAmount := basePrice - discount

	// Create subscription
	subscription := &model.Subscription{
		UserId:    userId,
		Plan:      plan,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 1, 0), // 1 month subscription
		IsActive:  true,
	}

	err := s.subscriptionRepo.CreateSubscription(subscription)
	if err != nil {
		return nil, errors.New("failed to create subscription")
	}

	// Create purchase record
	purchase := &model.Purchase{
		UserId:         userId,
		SubscriptionId: subscription.Id,
		Amount:         basePrice,
		Discount:       discount,
		Total:          totalAmount,
		VoucherCode:    voucherCode,
		PurchaseDate:   time.Now(),
	}

	err = s.purchaseRepo.CreatePurchase(purchase)
	if err != nil {
		return nil, errors.New("failed to create purchase")
	}

	return purchase, nil
}
