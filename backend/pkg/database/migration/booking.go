package migration

import (
	"context"
	"main/config"
	"main/pkg/database"

	"go.mongodb.org/mongo-driver/mongo"
)

func bookingDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("booking_db")
}