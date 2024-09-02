package repository

import (
	"context"
	"main/modules/booking/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository interface {
    CreateBooking(ctx context.Context, booking *models.Booking) error
    GetBookingByID(ctx context.Context, id string) (*models.Booking, error)
    ListBookingsByUser(ctx context.Context, userID string) ([]*models.Booking, error)
    UpdateBookingStatus(ctx context.Context, id, status string) error
}

type MongoBookingRepository struct {
    collection *mongo.Collection
}

func NewMongoBookingRepository(db *mongo.Database) BookingRepository {
    return &MongoBookingRepository{
        collection: db.Collection("bookings"),
    }
}

func (r *MongoBookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
    booking.ID = primitive.NewObjectID()
    booking.CreatedAt = primitive.NewDateTimeFromTime(time.Now()).Time().Unix()
    booking.UpdatedAt = booking.CreatedAt
    _, err := r.collection.InsertOne(ctx, booking)
    return err
}

func (r *MongoBookingRepository) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    var booking models.Booking
    err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
    return &booking, err
}

func (r *MongoBookingRepository) ListBookingsByUser(ctx context.Context, userID string) ([]*models.Booking, error) {
    cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    var bookings []*models.Booking
    err = cursor.All(ctx, &bookings)
    return bookings, err
}

func (r *MongoBookingRepository) UpdateBookingStatus(ctx context.Context, id, status string) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": status, "updated_at": primitive.NewDateTimeFromTime(time.Now()).Time().Unix()}})
    return err
}
