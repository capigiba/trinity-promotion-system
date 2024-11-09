package campaign

// CreateCampaignRequest represents the request payload for creating a campaign
type CreateCampaignRequest struct {
	Name        string  `json:"name" binding:"required"`
	Discount    float64 `json:"discount" binding:"required,gt=0"`
	MaxUsers    int     `json:"max_users" binding:"required,gt=0"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

// GenerateVouchersRequest represents the request payload for generating vouchers
type GenerateVouchersRequest struct {
	Count int `json:"count" binding:"required,gt=0"`
}
