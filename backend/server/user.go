package server

import (
	"main/modules/user/handlers"
	"main/modules/user/repository"
	"main/modules/user/usecase"
)
func (s *server) userService() {
	repo := repository.NewUserRepository(s.db)
	userUsecase := usecase.NewUserUsecase(repo)
	httpHandler := handlers.NewUserHttpHandler(s.cfg, userUsecase)
	user := s.app.Group("/user_v1")

	// Health check endpoint
	// user.GET("", s.healthCheckService)

	// User endpoints
	user.POST("/users/register", httpHandler.CreateUser())
	user.GET("/users/:user_id", httpHandler.FindOneUserProfile())
}

