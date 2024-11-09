package model

import "time"

type Campaign struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Discount    float64   `bson:"discount" json:"discount"`
	MaxUsers    int       `bson:"max_users" json:"max_users"`
	UsedUsers   int       `bson:"used_users" json:"used_users"`
	StartDate   time.Time `bson:"start_date" json:"start_date"`
	EndDate     time.Time `bson:"end_date" json:"end_date"`
	Description string    `bson:"description" json:"description"`
}
