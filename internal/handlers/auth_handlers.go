package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/ktariayman/go-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := c.BodyParser(&data); err != nil {
			log.Println("Error parsing request body:", err)
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		}

		user := models.User{}
		if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "invalid credentials"})
		}

		log.Println("Stored hashed password:", user.Password)
		log.Println("Provided password:", data.Password)

		// Compare passwords directly
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
			log.Println("Password comparison failed:", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "invalid credentials"})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": user.ID,
			"exp":    time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "could not login"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"token": tokenString})
	}
}

func LogoutUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "user logged out successfully"})
	}
}
