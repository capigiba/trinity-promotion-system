package model

import "time"

type Voucher struct {
	Id         string    `bson:"_id,omitempty" json:"id"`
	Code       string    `bson:"code" json:"code"`
	CampaignID string    `bson:"campaign_id" json:"campaign_id"`
	UserId     string    `bson:"user_id,omitempty" json:"user_id"`
	Used       bool      `bson:"used" json:"used"`
	ExpiryDate time.Time `bson:"expiry_date" json:"expiry_date"`
}
