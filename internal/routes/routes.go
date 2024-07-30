package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ktariayman/go-api/internal/handlers"
	auth "github.com/ktariayman/go-api/internal/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	api.Post("/register", handlers.RegisterUser(db))
	api.Post("/login", handlers.LoginUser(db))
	api.Post("/logout", handlers.LogoutUser())

	api.Post("/event", auth.Protected(), handlers.CreateEvent(db))
	api.Delete("/event/:id", auth.Protected(), handlers.DeleteEvent(db))
	api.Put("/event/:id", auth.Protected(), handlers.UpdateEvent(db))

	api.Get("/event/:id", handlers.GetEventByID(db))
	api.Get("/event", handlers.GetEvents(db))
	
	api.Post("/event/:id/participate", auth.Protected(), handlers.ParticipateInEvent(db))
	api.Post("/event/:id/cancel", auth.Protected(), handlers.CancelParticipation(db))

	api.Get("/user", auth.Protected(), handlers.GetAllUsers(db))
	api.Delete("/user/:id", auth.Protected(), handlers.DeleteUser(db))
	api.Put("/user/:id", auth.Protected(), handlers.UpdateUser(db))
}
