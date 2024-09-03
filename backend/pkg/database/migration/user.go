package migration

import (
	"context"
	"log"
	"main/config"
	"main/modules/user"
	"main/pkg/database"
	"main/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func userDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("user_db")
}

func UserMigrate(pctx context.Context, cfg *config.Config) {

	db := userDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("user_transactions")

	//set indexing
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
	})
	log.Println(indexs)

	col = db.Collection("users")

	//set indexing
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "_id", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
	})
	log.Println(indexs)

	documents := func() []any {
		roles := []*user.User{
			{
				Email:    "6531503042@lamduan.mfu.ac.th",
				//Hashing
				Password: func() string {
					HashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(HashedPassword)
				}(),
				Name: "Nimit Tanboontor",
				UserRoles: []user.UserRole{
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:		"admin001@gmail.com",
				Password:	func() string {
					HashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(HashedPassword)
				}(),
				Name:		"Admin01",
				UserRoles:	[]user.UserRole{
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:		"user01@gmail.com",
				Password:	func() string {
					HashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
					return string(HashedPassword)
				}(),
				Name:		"User01",
				UserRoles:	[]user.UserRole{
					{
						RoleTitle: "user",
						RoleCode:  2,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
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
	log.Println("Migrate user completed: ",results)
}