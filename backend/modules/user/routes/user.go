package routes

import (
	"main/modules/user/usecase"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App, uc usecase.IUserUsecase) {
	
		handler := &UserHandler{
			UserUsecase: uc,
		}
	
		app.Post("/users", handler.CreateUser)
		app.Get("/users/:id", handler.GetUserByID)
		app.Get("/users", handler.GetAllUsers)
		app.Put("/users/:id", handler.UpdateUser)
		app.Delete("/users/:id", handler.DeleteUser)
	}