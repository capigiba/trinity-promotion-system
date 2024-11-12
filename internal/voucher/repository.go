package voucher

import (
	"context"
	"time"
	"trinity/internal/model"

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
}

// NewRepository creates a new Voucher repository
func NewRepository(db *mongo.Database) Repository {
	return &repository{
		collection: db.Collection("vouchers"),
	}
}

// CreateVoucher inserts a new voucher into the database
func (r *repository) CreateVoucher(voucher *model.Voucher) error {
	_, err := r.collection.InsertOne(context.Background(), voucher)
	return err
}

// GetVoucherByCode retrieves a voucher by its code
func (r *repository) GetVoucherByCode(code string) (*model.Voucher, error) {
	var voucher model.Voucher
	err := r.collection.FindOne(context.Background(), map[string]interface{}{"code": code}).Decode(&voucher)
	return &voucher, err
}

// UpdateVoucher updates an existing voucher in the database
func (r *repository) UpdateVoucher(voucher *model.Voucher) error {
	voucher.Used = true
	_, err := r.collection.UpdateOne(context.Background(), map[string]interface{}{"_id": voucher.Id},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"used":    voucher.Used,
				"user_id": voucher.UserId,
				"updated": time.Now(),
			},
		})
	return err
}
