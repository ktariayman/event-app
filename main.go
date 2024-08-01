package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/ktariayman/go-api/docs"
	"github.com/ktariayman/go-api/internal/models"
	"github.com/ktariayman/go-api/internal/routes"
	"github.com/ktariayman/go-api/internal/seed"
	"github.com/ktariayman/go-api/pkg/config"
	"github.com/ktariayman/go-api/pkg/database"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	err = models.MigrateEvents(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	seed.Seed(db)    
	
	app := fiber.New()    
    app.Static("/openapi.yaml", "openapi.yaml")
    app.Get("/swagger/*", swagger.New(swagger.Config{
        URL: "/openapi.yaml",
    }))
	routes.SetupRoutes(app, db)
	log.Fatal(app.Listen(":8080"))
}
