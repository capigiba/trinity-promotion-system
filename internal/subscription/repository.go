package subscription

import (
	"context"
	"trinity/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateSubscription(subscription *model.Subscription) error
	// Additional methods if needed
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("subscriptions"),
	}
}

func (r *repository) CreateSubscription(subscription *model.Subscription) error {
	_, err := r.collection.InsertOne(context.Background(), subscription)
	return err
}
