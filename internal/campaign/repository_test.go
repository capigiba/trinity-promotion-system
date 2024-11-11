package campaign

import (
	"context"
	"testing"
	"time"

	"trinity/internal/model"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getTestDB connects to the test MongoDB instance
func getTestDB(t *testing.T) *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	assert.NoError(t, err, "Failed to connect to MongoDB")

	// Ping the database to verify connection
	err = client.Ping(context.Background(), nil)
	assert.NoError(t, err, "Failed to ping MongoDB")

	// Use a separate database for testing
	db := client.Database("trinity_test")

	// Clean up the database before and after tests
	t.Cleanup(func() {
		err := db.Drop(context.Background())
		assert.NoError(t, err, "Failed to drop test database")
		err = client.Disconnect(context.Background())
		assert.NoError(t, err, "Failed to disconnect MongoDB client")
	})

	return db
}

func TestRepository_CreateCampaign(t *testing.T) {
	db := getTestDB(t)
	repoInterface := NewRepository(db)

	// Perform type assertion to access the concrete *repository type
	repo, ok := repoInterface.(*repository)
	assert.True(t, ok, "Repository interface does not hold a *repository type")

	campaign := &model.Campaign{
		Name:        "Integration Test Campaign",
		Discount:    15,
		MaxUsers:    50,
		UsedUsers:   0,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(72 * time.Hour),
		Description: "A campaign for integration testing",
	}

	id, err := repo.CreateCampaign(campaign)
	assert.NoError(t, err, "CreateCampaign should not return an error")
	assert.NotEmpty(t, id, "Campaign ID should not be empty")
	assert.Equal(t, id, campaign.Id, "Returned ID should match campaign's ID")

	// Verify that the campaign was inserted into the database
	var insertedCampaign model.Campaign
	objID, err := primitive.ObjectIDFromHex(id)
	assert.NoError(t, err, "Invalid ObjectID format")

	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&insertedCampaign)
	assert.NoError(t, err, "Failed to find inserted campaign")
	assert.Equal(t, campaign.Name, insertedCampaign.Name, "Campaign name should match")
	assert.Equal(t, campaign.Discount, insertedCampaign.Discount, "Campaign discount should match")
}

func TestRepository_GetCampaignByID(t *testing.T) {
	db := getTestDB(t)
	repoInterface := NewRepository(db)

	// Perform type assertion to access the concrete *repository type
	repo, ok := repoInterface.(*repository)
	assert.True(t, ok, "Repository interface does not hold a *repository type")

	// First, insert a campaign to retrieve
	campaign := &model.Campaign{
		Name:        "Retrieve Test Campaign",
		Discount:    20,
		MaxUsers:    100,
		UsedUsers:   10,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(48 * time.Hour),
		Description: "A campaign for retrieval testing",
	}

	id, err := repo.CreateCampaign(campaign)
	assert.NoError(t, err, "CreateCampaign should not return an error")
	assert.NotEmpty(t, id, "Campaign ID should not be empty")

	// Test retrieving the campaign by ID
	retrievedCampaign, err := repo.GetCampaignByID(id)
	assert.NoError(t, err, "GetCampaignByID should not return an error")
	assert.Equal(t, campaign.Name, retrievedCampaign.Name, "Campaign name should match")
	assert.Equal(t, campaign.Discount, retrievedCampaign.Discount, "Campaign discount should match")

	// Test retrieving a non-existent campaign
	_, err = repo.GetCampaignByID("nonexistentid123")
	assert.Error(t, err, "GetCampaignByID should return an error for invalid ID")
}

// func TestRepository_IncrementUsedUsers(t *testing.T) {
// 	db := getTestDB(t)
// 	repoInterface := NewRepository(db)

// 	// Perform type assertion to access the concrete *repository type
// 	repo, ok := repoInterface.(*repository)
// 	assert.True(t, ok, "Repository interface does not hold a *repository type")

// 	// Insert a campaign to update
// 	campaign := &model.Campaign{
// 		Name:        "Increment Test Campaign",
// 		Discount:    10,
// 		MaxUsers:    100,
// 		UsedUsers:   25,
// 		StartDate:   time.Now(),
// 		EndDate:     time.Now().Add(24 * time.Hour),
// 		Description: "A campaign for increment testing",
// 	}

// 	id, err := repo.CreateCampaign(campaign)
// 	assert.NoError(t, err, "CreateCampaign should not return an error")

// 	// Increment used_users
// 	err = repo.IncrementUsedUsers(id)
// 	assert.NoError(t, err, "IncrementUsedUsers should not return an error")

// 	// Verify that used_users has been incremented
// 	updatedCampaign, err := repo.GetCampaignByID(id)
// 	assert.NoError(t, err, "GetCampaignByID should not return an error")
// 	assert.Equal(t, campaign.UsedUsers+1, updatedCampaign.UsedUsers, "UsedUsers should be incremented by 1")

// 	// Test incrementing a non-existent campaign
// 	err = repo.IncrementUsedUsers("nonexistentid123")
// 	assert.Error(t, err, "IncrementUsedUsers should return an error for invalid ID")
// }

func TestRepository_ListCampaigns(t *testing.T) {
	db := getTestDB(t)
	repoInterface := NewRepository(db)

	// Perform type assertion to access the concrete *repository type
	repo, ok := repoInterface.(*repository)
	assert.True(t, ok, "Repository interface does not hold a *repository type")

	// Insert multiple campaigns
	campaigns := []model.Campaign{
		{
			Name:        "List Test Campaign 1",
			Discount:    5,
			MaxUsers:    50,
			UsedUsers:   10,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(24 * time.Hour),
			Description: "First campaign for listing",
		},
		{
			Name:        "List Test Campaign 2",
			Discount:    10,
			MaxUsers:    100,
			UsedUsers:   20,
			StartDate:   time.Now(),
			EndDate:     time.Now().Add(48 * time.Hour),
			Description: "Second campaign for listing",
		},
	}

	for _, c := range campaigns {
		_, err := repo.CreateCampaign(&c)
		assert.NoError(t, err, "CreateCampaign should not return an error")
	}

	// Retrieve all campaigns
	retrievedCampaigns, err := repo.ListCampaigns()
	assert.NoError(t, err, "ListCampaigns should not return an error")
	assert.Len(t, retrievedCampaigns, len(campaigns), "Number of retrieved campaigns should match inserted campaigns")

	// Optionally, verify the contents
	for _, inserted := range campaigns {
		found := false
		for _, retrieved := range retrievedCampaigns {
			if retrieved.Name == inserted.Name {
				found = true
				assert.Equal(t, inserted.Discount, retrieved.Discount, "Discount should match")
				assert.Equal(t, inserted.MaxUsers, retrieved.MaxUsers, "MaxUsers should match")
				assert.Equal(t, inserted.UsedUsers, retrieved.UsedUsers, "UsedUsers should match")
				break
			}
		}
		assert.True(t, found, "Inserted campaign should be found in retrieved campaigns")
	}
}
