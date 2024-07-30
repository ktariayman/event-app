package seed

import (
	"log"
	"time"

	"github.com/ktariayman/go-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: string(adminPassword),
	}
	db.Create(&admin)

	userPassword, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	user := models.User{
		Name:     "Regular User",
		Email:    "user@example.com",
		Password: string(userPassword),
	}
	db.Create(&user)

	event := models.Event{
		Title:       "Sample Event",
		Description: "This is a sample event",
		Date:        time.Now().Format("2006-01-02"),
		Location:    "Sample Location",
		UserID:      admin.ID,
	}
	db.Create(&event)
}
