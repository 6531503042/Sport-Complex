package repository

import (
    "context"
    "errors"
    "fmt"
    "log"
    "main/modules/payment"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepositoryService interface {
    InsertPayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error)
    UpdatePayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error)
    FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error)
    FindPaymentsByUser(ctx context.Context, userId string) ([]payment.PaymentEntity, error)
    FindSlipByUserId(ctx context.Context, userId string) ([]payment.PaymentSlip, error)
    SaveSlip(ctx context.Context, slip payment.PaymentSlip) error
    UpdateSlipStatus(ctx context.Context, slipId string, newStatus string) error
    GetPendingSlips (ctx context.Context) ([]payment.PaymentSlip, error) 
}

type paymentRepository struct {
    db *mongo.Client
}

func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
    return &paymentRepository{db: db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
    return r.db.Database("payment_db")
}

func (r *paymentRepository) bookingDbConn(ctx context.Context) *mongo.Database {
    return r.db.Database("booking_db")
}

func (r *paymentRepository) slipDbConn(pctx context.Context) *mongo.Database {
    return r.db.Database("slip_db")
}


func (r *paymentRepository) InsertPayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error) {
    db := r.db.Database("payment_db")
    col := db.Collection("payments")

    // Set created and updated timestamps
    payment.CreatedAt = time.Now()
    payment.UpdatedAt = time.Now()

    // Insert payment into payments collection
    if _, err := col.InsertOne(ctx, payment); err != nil {
        log.Printf("Error: InsertPayment failed: %s", err.Error())
        return nil, fmt.Errorf("InsertPayment failed: %w", err)
    }

    return payment, nil
}

func (r *paymentRepository) FindPayment(ctx context.Context, paymentId string) (*payment.PaymentEntity, error) {
    db := r.db.Database("payment_db")
    col := db.Collection("payments")

    // Convert paymentId string to ObjectID
    objectId, err := primitive.ObjectIDFromHex(paymentId)
    if err != nil {
        log.Printf("Error: Invalid ObjectID: %s", err.Error())
        return nil, fmt.Errorf("error: invalid payment ID format")
    }

    // Find the payment by ObjectID
    result := new(payment.PaymentEntity)
    err = col.FindOne(ctx, bson.M{"_id": objectId}).Decode(result)
    if err != nil {
        log.Printf("Error: FindPayment failed: %s", err.Error())
        return nil, fmt.Errorf("error: FindPayment failed")
    }

    return result, nil
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *payment.PaymentEntity) (*payment.PaymentEntity, error) {
    db := r.db.Database("payment_db")
    col := db.Collection("payments")

    _, err := col.UpdateOne(ctx, bson.M{"_id": payment.Id}, bson.M{"$set": payment})
    if err != nil {
        log.Printf("Error: UpdatePayment failed: %s", err.Error())
        return nil, fmt.Errorf("error: UpdatePayment failed")
    }

    return payment, nil
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

func (r *paymentRepository)  SaveSlip(ctx context.Context, slip payment.PaymentSlip) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    db := r.slipDbConn(ctx)
    col := db.Collection("slips")

    _, err := col.InsertOne(ctx, slip)
    if err != nil {
        log.Printf("Error: SaveSlip failed: %s", err.Error())
        return errors.New("error: SaveSlip failed")
    }

    return nil
}

func (r *paymentRepository) FindSlipByUserId(ctx context.Context, userId string) ([]payment.PaymentSlip, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Establish a connection to the database and collection
    db := r.slipDbConn(ctx)
    col := db.Collection("slips")

    // Find slips matching the user_id
    cursor, err := col.Find(ctx, bson.M{"user_id": userId})
    if err != nil {
        log.Printf("Error: FindSlipByUserId failed to execute query: %s", err.Error())
        return nil, errors.New("error: failed to retrieve slips by user ID")
    }
    defer cursor.Close(ctx) // Ensure cursor is closed after use

    // Decode the results directly into a slice of PaymentSlip
    var result []payment.PaymentSlip
    if err := cursor.All(ctx, &result); err != nil {
        log.Printf("Error: FindSlipByUserId failed to decode results: %s", err.Error())
        return nil, errors.New("error: failed to decode slips for user ID")
    }

    return result, nil
}

func (r *paymentRepository) UpdateSlipStatus(ctx context.Context, slipId string, newStatus string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Update status in payment_db's slips collection
    paymentDb := r.slipDbConn(ctx)
    paymentCol := paymentDb.Collection("slips")

    _, err := paymentCol.UpdateOne(ctx, bson.M{"_id": slipId}, bson.M{"$set": bson.M{"status": newStatus}})
    if err != nil {
        log.Printf("Error: UpdateSlipStatus failed to update slip in payment_db: %s", err.Error())
        return errors.New("error: failed to update slip status in payment database")
    }

    // Fetch the booking ID associated with this slip
    var slip payment.PaymentSlip
    if err := paymentCol.FindOne(ctx, bson.M{"_id": slipId}).Decode(&slip); err != nil {
        log.Printf("Error: UpdateSlipStatus failed to fetch slip: %s", err.Error())
        return errors.New("error: failed to fetch slip details")
    }

    // Update status in booking_db's booking_transaction collection
    bookingDb := r.bookingDbConn(ctx) // Connect to booking database
    bookingCol := bookingDb.Collection("booking_transaction")

    _, err = bookingCol.UpdateOne(ctx, bson.M{"booking_id": slip.BookingID}, bson.M{"$set": bson.M{"status": newStatus}})
    if err != nil {
        log.Printf("Error: UpdateSlipStatus failed to update booking transaction in booking_db: %s", err.Error())
        return errors.New("error: failed to update booking transaction status in booking database")
    }

    return nil
}

func (r *paymentRepository) GetPendingSlips(ctx context.Context) ([]payment.PaymentSlip, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    // Fetch pending slips from the payment_db
    db := r.slipDbConn(ctx)
    col := db.Collection("slips")

    cursor, err := col.Find(ctx, bson.M{"status": "pending"})
    if err != nil {
        log.Printf("Error: GetPendingSlips failed to execute query: %s", err.Error())
        return nil, errors.New("error: failed to retrieve pending slips")
    }
    defer cursor.Close(ctx)

    // Decode the results into a slice of PaymentSlip
    var pendingSlips []payment.PaymentSlip
    if err := cursor.All(ctx, &pendingSlips); err != nil {
        log.Printf("Error: GetPendingSlips failed to decode results: %s", err.Error())
        return nil, errors.New("error: failed to decode pending slips")
    }

    return pendingSlips, nil
}
