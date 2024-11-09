package campaign

import (
	"context"

	"trinity/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines campaign data access methods
type Repository interface {
	CreateCampaign(campaign *model.Campaign) error
	GetCampaignByID(id string) (*model.Campaign, error)
	IncrementUsedUsers(id string) error
}

// repository implements Repository interface
type repository struct {
	collection *mongo.Collection
}

// NewRepository creates a new Campaign repository
func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("campaigns"),
	}
}

// CreateCampaign inserts a new campaign into the database
func (r *repository) CreateCampaign(campaign *model.Campaign) error {
	_, err := r.collection.InsertOne(context.Background(), campaign)
	return err
}

// GetCampaignByID retrieves a campaign by its ID
func (r *repository) GetCampaignByID(id string) (*model.Campaign, error) {
	var campaign model.Campaign
	err := r.collection.FindOne(context.Background(), map[string]interface{}{"_id": id}).Decode(&campaign)
	return &campaign, err
}

// IncrementUsedUsers increments the used_users field of a campaign
func (r *repository) IncrementUsedUsers(id string) error {
	_, err := r.collection.UpdateOne(context.Background(), map[string]interface{}{"_id": id},
		map[string]interface{}{
			"$inc": map[string]interface{}{"used_users": 1},
		})
	return err
}
