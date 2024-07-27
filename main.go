package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/ktariayman/go-api/models"
	"github.com/ktariayman/go-api/storage"
	"gorm.io/gorm"
)
type Repo struct {
	DB *gorm.DB 
}
type Event struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
}



func (r * Repo) CreateEvent (context *fiber.Ctx) error{
	event := Event{}
	err := context.BodyParser(&event)	

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create event"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "event has been added"})
	return nil
}

func (r *Repo) DeleteEvent(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(&models.Event{}, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete event"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "event deleted successfully"})
	return nil
}

func (r *Repo) GetEvents(context *fiber.Ctx) error {
	events := &[]models.Event{}

	err := r.DB.Find(events).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get events"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "events fetched successfully", "data": events})
	return nil
}

func (r *Repo) GetEventByID(context *fiber.Ctx) error {
	id := context.Params("id")
	event := &models.Event{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get the event"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "event fetched successfully", "data": event})
	return nil
}

func (r *Repo) UpdateEvent(context *fiber.Ctx) error {
	id := context.Params("id")
	event := &models.Event{}

	err := r.DB.Where("id = ?", id).First(event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "event not found"})
		return err
	}

	err = context.BodyParser(event)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Save(event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update event"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "event updated successfully"})
	return nil
}
func (r *Repo) SetupRoutes(app *fiber.App) {
api := app.Group("/api")
	api.Post("/create_event", r.CreateEvent)
	api.Delete("/delete_event/:id", r.DeleteEvent)
	api.Get("/get_events/:id", r.GetEventByID)
	api.Get("/events", r.GetEvents)
	api.Put("/update_event/:id", r.UpdateEvent)
}
func main(){
	err := godotenv.Load(".env")

		if err != nil {
			log.Fatal(err)
		}
		config := &storage.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Password: os.Getenv("DB_PASS"),
			User:     os.Getenv("DB_USER"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			DBName:   os.Getenv("DB_NAME"),
		}
		db, err := storage.NewConnection(config)

		if err != nil {
			log.Fatal("could not load the database")
		}
		err = models.MigrateEvents(db)
		if err != nil {
			log.Fatal("could not migrate db")
		}

		r:= Repo{
			DB: db,
		}
		app:= fiber.New()
		r.SetupRoutes(app)
		app.Listen(":8080")
	}