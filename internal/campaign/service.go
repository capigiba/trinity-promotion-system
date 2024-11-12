package campaign

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"trinity/internal/model"
	"trinity/internal/voucher"
	"trinity/pkg/logger"
)

// Service defines campaign business logic methods
type Service interface {
	CreateCampaign(campaign *model.Campaign) (string, error)
	GenerateVouchers(campaignID string, count int) ([]model.Voucher, error)
	ListCampaigns() ([]model.Campaign, error)
}

// service implements Service interface
type service struct {
	repo        Repository
	voucherRepo voucher.Repository
	logger      logger.Logger
}

// NewService creates a new Campaign service
func NewService(repo Repository, voucherRepo voucher.Repository) Service {
	return &service{
		repo:        repo,
		voucherRepo: voucherRepo,
		logger:      logger.NewLogger("campaignService"),
	}
}

// CreateCampaign creates a new campaign
func (s *service) CreateCampaign(campaign *model.Campaign) (string, error) {
	// Validate campaign dates
	if campaign.StartDate.After(campaign.EndDate) {
		return "", errors.New("start date must be before end date")
	}

	// Create campaign
	id, err := s.repo.CreateCampaign(campaign)
	if err != nil {
		s.logger.Errorf("Failed to create campaign: %v", err)
		return "", err
	}

	return id, nil
}

// GenerateVouchers generates vouchers for a campaign
func (s *service) GenerateVouchers(campaignID string, count int) ([]model.Voucher, error) {
	campaign, err := s.repo.GetCampaignByID(campaignID)
	if err != nil {
		s.logger.Errorf("Failed to get campaign: %v", err)
		return nil, err
	}

	remainingVouchers := campaign.MaxUsers - campaign.UsedUsers
	if count > remainingVouchers {
		return nil, fmt.Errorf("not enough vouchers remaining: requested %d, available %d", count, remainingVouchers)
	}

	var generatedVouchers []model.Voucher

	for i := 0; i < count; i++ {
		code := s.generateVoucherCode()
		if code == "" {
			continue
		}
		voucher := model.Voucher{
			Code:       code,
			CampaignID: campaignID,
			Used:       false,
			ExpiryDate: campaign.EndDate,
		}
		err := s.voucherRepo.CreateVoucher(&voucher)
		if err != nil {
			s.logger.Errorf("Failed to create voucher: %v", err)
			continue
		}
		generatedVouchers = append(generatedVouchers, voucher)
	}

	// Increment used users
	err = s.repo.IncrementUsedUsers(campaignID, len(generatedVouchers))
	if err != nil {
		s.logger.Errorf("Failed to increment used users: %v", err)
		return nil, err
	}

	return generatedVouchers, nil
}

// generateVoucherCode generates a unique voucher code
func (s *service) generateVoucherCode() string {
	bytes := make([]byte, 5) // 10 characters when hex-encoded
	if _, err := rand.Read(bytes); err != nil {
		s.logger.Errorf("Failed to generate random bytes: %v", err)
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// ListCampaigns retrieves all campaigns
func (s *service) ListCampaigns() ([]model.Campaign, error) {
	return s.repo.ListCampaigns()
}
