package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Booking struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"`
    UserID     string             `bson:"user_id"`
    FacilityID string             `bson:"facility_id"`
    Timeslot   string             `bson:"timeslot"`
    Price      float64            `bson:"price"`
    Status     string             `bson:"status"`
    CreatedAt  int64              `bson:"created_at"`
    UpdatedAt  int64              `bson:"updated_at"`
}
