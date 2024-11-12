package voucher

import (
	"context"
	"errors"
	"fmt"
	"time"
	"trinity/internal/model"
	"trinity/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines voucher data access methods
type Repository interface {
	CreateVoucher(voucher *model.Voucher) error
	GetVoucherByCode(code string) (*model.Voucher, error)
	UpdateVoucher(voucher *model.Voucher) error
}

// repository implements Repository interface
type repository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

// NewRepository creates a new Voucher repository
func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("vouchers"),
		logger:     logger.NewLogger("voucherRepository"),
	}
}

// CreateVoucher inserts a new voucher into the database
func (r *repository) CreateVoucher(voucher *model.Voucher) error {
	result, err := r.collection.InsertOne(context.Background(), voucher)
	if err != nil {
		r.logger.Errorf("Failed to insert voucher: %v", err)
		return err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Errorf("InsertedID is not an ObjectID")
		return fmt.Errorf("failed to convert InsertedID to ObjectID")
	}

	voucher.Id = oid.Hex()
	r.logger.Infof("Voucher created with ID: %s", voucher.Id)
	return nil
}

// GetVoucherByCode retrieves a voucher by its code
func (r *repository) GetVoucherByCode(code string) (*model.Voucher, error) {
	var voucher model.Voucher
	err := r.collection.FindOne(context.Background(), bson.M{"code": code}).Decode(&voucher)
	if err != nil {
		r.logger.Errorf("Failed to find voucher by code %s: %v", code, err)
		return nil, err
	}
	return &voucher, nil
}

// UpdateVoucher updates an existing voucher in the database
func (r *repository) UpdateVoucher(voucher *model.Voucher) error {
	objectId, err := primitive.ObjectIDFromHex(voucher.Id)
	if err != nil {
		r.logger.Errorf("Invalid voucher ID format: %v", err)
		return fmt.Errorf("invalid voucher ID: %v", err)
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"used":    voucher.Used,
			"user_id": voucher.UserId,
			"updated": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		r.logger.Errorf("Failed to update voucher: %v", err)
		return fmt.Errorf("failed to update voucher: %v", err)
	}

	if result.MatchedCount == 0 {
		r.logger.Errorf("No voucher found with ID: %s", voucher.Id)
		return errors.New("no voucher found with the provided ID")
	}

	if result.ModifiedCount == 0 {
		r.logger.Warnf("Voucher %s was not updated, possibly already marked as used", voucher.Id)
		return errors.New("voucher was not updated, possibly already marked as used")
	}

	r.logger.Infof("Voucher %s updated successfully. Matched: %d, Modified: %d",
		voucher.Id, result.MatchedCount, result.ModifiedCount)

	return nil
}
