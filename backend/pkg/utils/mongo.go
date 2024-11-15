package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConvertToObjectId(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return objectId
}

func NewMongoDBConnection(pctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
        ApplyURI(uri).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    return client, nil
}