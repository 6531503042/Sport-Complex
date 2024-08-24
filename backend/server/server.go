package server

import (
	"context"
	"main/config"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app		*echo.Echo
		db		*mongo.Client
		cfg		*config.Config
	}
)

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		app:	echo.New(),
		db:		db,
		cfg: 	cfg,
	}

	//Cors
	// s.app.Use(middleware.CORS


}

// func (s *server) setupMiddleware() {
// 	// Request Timeout
// 	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
// 		Skipper:      middleware.DefaultSkipper,
// 		ErrorMessage: "Error: Request Timeout",
// 		Timeout:      30 * time.Second,
// 	}))

// 	// CORS
// 	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
// 		Skipper:      middleware.DefaultSkipper,
// 		AllowOrigins: []string{"*"},
// 		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
// 	}))

// 	// Body Limit
// 	s.app.Use(middleware.BodyLimit("10M"))

// 	// Additional middleware can be configured here...
// }