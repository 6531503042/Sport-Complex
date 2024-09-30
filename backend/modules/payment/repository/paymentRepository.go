package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"main/modules/payment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	PaymentRepositoryService interface {
		InsertPayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error)
		UpdatePayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error)
		FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error)
		FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error)
	}

	paymentRepository struct {
		db *mongo.Client
	}
)

// NewPaymentRepository returns a new instance of PaymentRepositoryService using the given mongo client.
func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}

// InsertPayment inserts a new payment into the database
func (r *paymentRepository) InsertPayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payments")

	// Set created and updated timestamps
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	// Insert payment into payments collection
	_, err := col.InsertOne(ctx, payment)
	if err != nil {
		log.Printf("Error: InsertPayment failed: %s", err.Error())
		return nil, errors.New("error: InsertPayment failed")
	}

	return payment, nil
}

// UpdatePayment updates an existing payment in the database
func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payments")

	// Update the payment document
	_, err := col.UpdateOne(ctx, bson.M{"_id": payment.Id}, bson.M{"$set": payment})
	if err != nil {
		log.Printf("Error: UpdatePayment failed: %s", err.Error())
		return nil, errors.New("error: UpdatePayment failed")
	}

	return payment, nil
}

// FindPayment retrieves a payment by its ID
func (r *paymentRepository) FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payments")
	result := new(payment.PaymentEntity)

	err := col.FindOne(ctx, bson.M{"_id": paymentId}).Decode(result)
	if err != nil {
		log.Printf("Error: FindPayment failed: %s", err.Error())
		return nil, errors.New("error: FindPayment failed")
	}

	return result, nil
}

// FindPaymentsByUser retrieves all payments made by a specific user
func (r *paymentRepository) FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payments")

	cursor, err := col.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		log.Printf("Error: FindPaymentsByUser failed: %s", err.Error())
		return nil, errors.New("error: FindPaymentsByUser failed")
	}
	defer cursor.Close(ctx)

	var result []payment.PaymentEntity
	if err = cursor.All(ctx, &result); err != nil {
		log.Printf("Error: FindPaymentsByUser failed: %s", err.Error())
		return nil, errors.New("error: FindPaymentsByUser failed")
	}

	return result, nil
}
