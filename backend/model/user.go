package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Email     string             `bson:"email"`
    Name      string             `bson:"name"`
    Phone     string             `bson:"phone"`
    Password  string             `bson:"password"`
    CreatedAt string             `bson:"created_at"`
}
