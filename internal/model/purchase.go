package model

import "time"

type Purchase struct {
	Id             string    `bson:"_id,omitempty" json:"id"`
	UserId         string    `bson:"user_id" json:"user_id"`
	SubscriptionId string    `bson:"subscription_id" json:"subscription_id"`
	Amount         float64   `bson:"amount" json:"amount"`
	Discount       float64   `bson:"discount" json:"discount"`
	Total          float64   `bson:"total" json:"total"`
	VoucherCode    string    `bson:"voucher_code,omitempty" json:"voucher_code"`
	PurchaseDate   time.Time `bson:"purchase_date" json:"purchase_date"`
}
