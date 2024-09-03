package repository

import (
	"context"
	"errors"
	"log"
	"main/modules/booking"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type(
	BookingRepositoryService interface {
		InsertOneBooking(ctx context.Context, req *booking.Booking) (primitive.ObjectID, error)
		IsSlotFree (ctx context.Context, slot booking.TimeSlot) (bool, error)
		GetAllBooking (ctx context.Context) ([]booking.Booking, error)
	}

	BookingRepository struct {
		db *mongo.Client
	}
)

func NewBookingRepository(db *mongo.Client) BookingRepositoryService {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) bookingDbConn(ctx context.Context) *mongo.Database {
	return r.db.Database("booking_db")
}

func (r *BookingRepository) InsertOneBooking(ctx context.Context, req *booking.Booking) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error InsertOneBooking: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: insert one booking failed")
	}
	
	bookingId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("error: insert one booking failed")
	}

	return bookingId, nil
}

func (r *BookingRepository) IsSlotFree (ctx context.Context, slot booking.TimeSlot) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking")

	//Query from db
	filter := bson.M{
		"slot.start": bson.M{"$gt": slot.Start},
		"slot.end": bson.M{"$lt": slot.End},
	}

	slotCount, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: failed to query filter")
		return false, err
	}

	if slotCount > 0 {
		return false, nil
	}
	return true, nil
}

func (r *BookingRepository) GetAllBooking (ctx context.Context) ([]booking.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db := r.bookingDbConn(ctx)
	col := db.Collection("booking")

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: GetAllBooking: %s", err.Error())
		return nil, err
	}

	defer cursor.Close(ctx)

	var booking []booking.Booking
	if err = cursor.All(ctx, &booking); err != nil {
		log.Printf("Error GetAllBookings: %s", err.Error())
		return nil, errors.New("error: failed to parse bookings")
	}

	return booking, nil
}