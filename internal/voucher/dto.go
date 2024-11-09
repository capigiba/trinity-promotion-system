package voucher

// RedeemVoucherRequest represents the request payload for redeeming a voucher
type RedeemVoucherRequest struct {
	Code   string `json:"code" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}
