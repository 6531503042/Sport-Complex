package migration

import (
	"context"
	"log"
	"main/config"
	"main/modules/auth"
	"main/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("auth_db")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {

	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("auth")

	//Auth Indexing
	indexs, err := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "refresh_token", Value: 1}}},
	})
	if err != nil {
		panic(err)
	}

	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	//Role
	col = db.Collection("roles")

	indexs, err = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "code", Value: 1}}},
	})
	if err != nil {
		panic(err)
	}

	for _, index := range indexs {
		log.Printf("Index :%s", index)
	}

	// roles data
	documents := func() []any {
		roles := []*auth.Role{
			{
				Title: "admin",
				Code:  1,
			},
			{
				Title: "insider",
				Code:  2,
			},
			{
				Title: "outsider",
				Code:  3,
			},
			{
				Title: "guest",
				Code:  4,
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, documents, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate auth completed: ", results)
}

