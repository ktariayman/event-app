package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/ktariayman/go-api/auth"
	"github.com/ktariayman/go-api/models"
	"github.com/ktariayman/go-api/storage"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Event struct {
	gorm.Model
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Date         string         `json:"date"`
	Location     string         `json:"location"`
	UserID       uint           `json:"user_id"`
	Participants []models.User  `gorm:"many2many:event_participants;" json:"participants"`
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func (r *Repo) CreateEvent(context *fiber.Ctx) error {
	userID := context.Locals("userID").(float64) 
	event := Event{
		UserID: uint(userID),
	}
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
	userID := context.Locals("userID").(float64)
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	event := Event{}
	err := r.DB.Where("id = ? AND user_id = ?", id, uint(userID)).First(&event).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not find event"})
		return err
	}

	err = r.DB.Delete(&event).Error
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
	userID := context.Locals("userID").(float64)
	id := context.Params("id")
	event := &models.Event{}
	err := r.DB.Where("id = ? AND user_id = ?", id, uint(userID)).First(event).Error
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

func (r *Repo) RegisterUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	// Check if email already exists
	existingUser := User{}
	err = r.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "email already in use"})
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "could not hash password"})
		return err
	}
	user.Password = string(hashedPassword)

	err = r.DB.Create(&user).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not register user"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "user registered successfully"})
	return nil
}


func (r *Repo) LoginUser(context *fiber.Ctx) error {
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := context.BodyParser(&data)
	if err != nil {
		log.Println("Error parsing request body:", err)
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	user := User{}
	err = r.DB.Where("email = ?", data.Email).First(&user).Error
	if err != nil {
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "invalid credentials"})
		return err
	}

	log.Println("Stored hashed password:", user.Password)
	log.Println("Provided password:", data.Password)

	// Compare passwords directly
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		log.Println("Password comparison failed:", err)
		context.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "invalid credentials"})
		return err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "could not login"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"token": tokenString})
	return nil
}

func (r *Repo) LogoutUser(context *fiber.Ctx) error {
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "user logged out successfully"})
	return nil
}

func (r *Repo) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", r.RegisterUser)
	api.Post("/login", r.LoginUser)
	api.Post("/logout", r.LogoutUser)

	api.Post("/event", auth.Protected(), r.CreateEvent)
	api.Delete("/event/:id", auth.Protected(), r.DeleteEvent)
	api.Put("/event/:id", auth.Protected(), r.UpdateEvent)
	api.Get("/events/:id", r.GetEventByID)
	api.Get("/events", r.GetEvents)
}

func main() {
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
	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	r := Repo{DB: db}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
