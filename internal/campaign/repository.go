package campaign

import (
	"context"
	"fmt"

	"trinity/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines campaign data access methods
type Repository interface {
	CreateCampaign(campaign *model.Campaign) (string, error)
	GetCampaignByID(id string) (*model.Campaign, error)
	IncrementUsedUsers(id string) error
	ListCampaigns() ([]model.Campaign, error)
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

// CreateCampaign inserts a new campaign into the database and returns its ID
func (r *repository) CreateCampaign(campaign *model.Campaign) (string, error) {
	result, err := r.collection.InsertOne(context.Background(), campaign)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert to ObjectID")
	}

	campaign.Id = oid.Hex()
	return campaign.Id, nil
}

// GetCampaignByID retrieves a campaign by its ID
func (r *repository) GetCampaignByID(campaignID string) (*model.Campaign, error) {
	objID, err := primitive.ObjectIDFromHex(campaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	var campaign model.Campaign
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&campaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no documents in result")
		}
		return nil, err
	}

	return &campaign, nil
}

// IncrementUsedUsers increments the used_users field of a campaign
func (r *repository) IncrementUsedUsers(id string) error {
	_, err := r.collection.UpdateOne(context.Background(), map[string]interface{}{"_id": id},
		map[string]interface{}{
			"$inc": map[string]interface{}{"used_users": 1},
		})
	return err
}

// ListCampaigns retrieves all campaigns from the database
func (r *repository) ListCampaigns() ([]model.Campaign, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campaigns []model.Campaign
	for cursor.Next(context.Background()) {
		var campaign model.Campaign
		if err := cursor.Decode(&campaign); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return campaigns, nil
}
