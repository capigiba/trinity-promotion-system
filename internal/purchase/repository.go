package purchase

import (
	"context"

	"trinity/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines purchase data access methods
type Repository interface {
	CreatePurchase(purchase *model.Purchase) error
}

// repository implements Repository interface
type repository struct {
	collection *mongo.Collection
}

// NewRepository creates a new Purchase repository
func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("purchases"),
	}
}

// CreatePurchase inserts a new purchase into the database
func (r *repository) CreatePurchase(purchase *model.Purchase) error {
	_, err := r.collection.InsertOne(context.Background(), purchase)
	return err
}
