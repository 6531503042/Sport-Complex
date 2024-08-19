package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	User struct {
		Id	primitive.ObjectID	`json:"_id" bson:"_id,omitempty"`
		Email string `json:"email" bson:"email"`
		Username string `json:"username bson:"username"`
		Password string `json:"password" bson:"password"`
		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
		RoleCode []UserRole  `bson:"role_code"`
	}

	UserRole struct {
		RoleTitle string `json:"role_title" bson:"role_title"`
		RoleCode int `json:"role_code" bson:"role_code"`

	}

	UserProfileBson struct {
		Id primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Email string `json:"email" bson:"email"`
		Username string `json:"username bson:"username"`
		Password string `json:"password" bson:"password"`
		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	}
)