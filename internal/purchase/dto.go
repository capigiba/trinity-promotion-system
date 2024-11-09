package purchase

import "trinity/internal/model"

// ProcessPurchaseRequest represents the request payload for processing a purchase
type ProcessPurchaseRequest struct {
	UserId      string                 `json:"user_id" binding:"required"`
	Plan        model.SubscriptionPlan `json:"plan" binding:"required"`
	VoucherCode string                 `json:"voucher_code,omitempty"`
}
