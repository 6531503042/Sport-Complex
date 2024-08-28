package server

import (
	"main/modules/user/handlers"
	"main/modules/user/repository"
	"main/modules/user/usecase"
)

func (s *server) userService() {
	// Initialize the repository, usecase, and HTTP handler
	repo := repository.NewUserRepository(s.db)
	userUsecase := usecase.NewUserUsecase(repo)
	httpHandler := handlers.NewUserHttpHandler(s.cfg, userUsecase)

	// Define a route group for user-related endpoints
	user := s.app.Group("/user_v1")

	// Health check endpoint (uncomment and implement if needed)
	// user.GET("", s.healthCheckService)

	// User endpoints
	user.POST("/users/register", httpHandler.CreateUser) 
	user.GET("/users/:user_id", httpHandler.FindOneUserProfile) 
}