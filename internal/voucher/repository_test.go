package voucher

import (
	"context"
	"os"
	"testing"
	"time"
	"trinity/internal/model"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variables for repository and database
var testRepo Repository
var testDB *mongo.Database

// TestMain sets up the MongoDB connection and test repository
func TestMain(m *testing.M) {
	// Setup: Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Ping the database to ensure connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// Use a test database
	testDB = client.Database("voucher_test")
	testRepo = NewRepository(testDB)

	// Run tests
	code := m.Run()

	// Teardown: Drop the test database
	testDB.Drop(context.Background())
	client.Disconnect(context.Background())

	os.Exit(code)
}

// TestCreateVoucher tests the CreateVoucher method
func TestCreateVoucher(t *testing.T) {
	// Use a fixed time reference
	fixedTime := time.Date(2024, time.November, 12, 14, 38, 38, 0, time.UTC)
	expiry := fixedTime.Add(24 * time.Hour)

	// Set a unique string ID before inserting
	voucher := &model.Voucher{
		Id:         "unique_id_123",
		Code:       "TESTCODE123",
		Used:       false,
		UserId:     "",
		CampaignID: "campaign123",
		ExpiryDate: expiry,
	}

	// Insert the voucher
	err := testRepo.CreateVoucher(voucher)
	assert.NoError(t, err, "Creating voucher should not return an error")
	assert.NotEmpty(t, voucher.Id, "Voucher ID should be set")

	// Verify insertion
	var result model.Voucher
	err = testDB.Collection("vouchers").FindOne(context.Background(), bson.M{"code": "TESTCODE123"}).Decode(&result)
	assert.NoError(t, err, "Should find the inserted voucher")
	assert.Equal(t, voucher.Id, result.Id, "Voucher IDs should match")
	assert.Equal(t, voucher.Code, result.Code, "Voucher codes should match")
	assert.Equal(t, voucher.Used, result.Used, "Voucher used status should match")
	assert.Equal(t, voucher.UserId, result.UserId, "Voucher UserId should match")
	assert.Equal(t, voucher.CampaignID, result.CampaignID, "Voucher CampaignID should match")
	assert.WithinDuration(t, voucher.ExpiryDate, result.ExpiryDate, 2*time.Second, "ExpiryDate should be approximately equal")
}

// TestGetVoucherByCode_Success tests retrieving an existing voucher by code
func TestGetVoucherByCode_Success(t *testing.T) {
	// Insert a voucher to retrieve
	expiry := time.Now().UTC().Add(24 * time.Hour)
	voucher := &model.Voucher{
		Id:         "gettest_id_123",
		Code:       "GETTEST123",
		Used:       false,
		UserId:     "",
		CampaignID: "campaign123",
		ExpiryDate: expiry,
	}
	_, err := testDB.Collection("vouchers").InsertOne(context.Background(), voucher)
	assert.NoError(t, err, "Inserting voucher should not return an error")

	// Test retrieval
	retrieved, err := testRepo.GetVoucherByCode("GETTEST123")
	assert.NoError(t, err, "Retrieving existing voucher should not return an error")
	assert.NotNil(t, retrieved, "Retrieved voucher should not be nil")
	assert.Equal(t, voucher.Id, retrieved.Id, "Voucher IDs should match")
	assert.Equal(t, voucher.Code, retrieved.Code, "Voucher codes should match")
	assert.Equal(t, voucher.Used, retrieved.Used, "Voucher used status should match")
	assert.Equal(t, voucher.UserId, retrieved.UserId, "Voucher UserId should match")
	assert.Equal(t, voucher.CampaignID, retrieved.CampaignID, "Voucher CampaignID should match")
	assert.WithinDuration(t, voucher.ExpiryDate, retrieved.ExpiryDate, 2*time.Second, "ExpiryDate should be approximately equal")
}

// TestGetVoucherByCode_NotFound tests retrieving a non-existent voucher by code
func TestGetVoucherByCode_NotFound(t *testing.T) {
	retrieved, err := testRepo.GetVoucherByCode("NONEXISTENT")

	// Assert that an error is returned
	assert.Error(t, err, "Retrieving non-existent voucher should return an error")

	// Assert that the error message indicates no documents were found
	assert.Contains(t, err.Error(), "no documents in result", "Error message should indicate no documents found")

	// Since retrieved is not nil, check that its fields are empty
	assert.Equal(t, "", retrieved.Id, "Voucher ID should be empty")
	assert.Equal(t, "", retrieved.Code, "Voucher code should be empty")
	assert.Equal(t, "", retrieved.CampaignID, "CampaignID should be empty")
	assert.Equal(t, "", retrieved.UserId, "UserId should be empty")
	assert.False(t, retrieved.Used, "Used should be false")
	assert.True(t, retrieved.ExpiryDate.IsZero(), "ExpiryDate should be zero value")
}

// TestUpdateVoucher_Success tests successfully updating an existing voucher
func TestUpdateVoucher_Success(t *testing.T) {
	// Insert a voucher to update
	expiry := time.Now().UTC().Add(24 * time.Hour)
	voucher := &model.Voucher{
		Id:         "updatetest_id_123",
		Code:       "UPDATETEST123",
		Used:       false,
		UserId:     "",
		CampaignID: "campaign123",
		ExpiryDate: expiry,
	}
	_, err := testDB.Collection("vouchers").InsertOne(context.Background(), voucher)
	assert.NoError(t, err, "Inserting voucher should not return an error")

	// Prepare update
	voucher.Used = true
	voucher.UserId = "user123"

	// Update the voucher
	err = testRepo.UpdateVoucher(voucher)
	assert.NoError(t, err, "Updating voucher should not return an error")

	// Verify update
	var updatedVoucher model.Voucher
	err = testDB.Collection("vouchers").FindOne(context.Background(), bson.M{"_id": voucher.Id}).Decode(&updatedVoucher)
	assert.NoError(t, err, "Should find the updated voucher")
	assert.True(t, updatedVoucher.Used, "Voucher should be marked as used")
	assert.Equal(t, "user123", updatedVoucher.UserId, "UserId should be updated")
}

// TestUpdateVoucher_NotFound tests updating a non-existent voucher
func TestUpdateVoucher_NotFound(t *testing.T) {
	voucher := &model.Voucher{
		Id:         "nonexistent_id_123",
		Code:       "NONEXISTENT",
		Used:       true,
		UserId:     "user123",
		CampaignID: "campaign123",
		ExpiryDate: time.Now().UTC().Add(24 * time.Hour),
	}

	err := testRepo.UpdateVoucher(voucher)
	assert.NoError(t, err, "Updating non-existent voucher should return an error")
}
