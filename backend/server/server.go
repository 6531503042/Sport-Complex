package server

import (
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app		*echo.Echo
		db		*mongo.Client
	}
)