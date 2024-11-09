package model

import "time"

type SubscriptionPlan string

const (
	PlanSilver SubscriptionPlan = "silver"
	PlanGold   SubscriptionPlan = "gold"
)

type Subscription struct {
	Id        string           `bson:"_id,omitempty" json:"id"`
	UserId    string           `bson:"user_id" json:"user_id"`
	Plan      SubscriptionPlan `bson:"plan" json:"plan"`
	StartDate time.Time        `bson:"start_date" json:"start_date"`
	EndDate   time.Time        `bson:"end_date" json:"end_date"`
	IsActive  bool             `bson:"is_active" json:"is_active"`
}
