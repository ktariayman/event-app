package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ktariayman/go-api/internal/helpers"
	"github.com/ktariayman/go-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := models.User{}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		}

		existingUser := models.User{}
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "email already in use"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "could not hash password"})
		}
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "could not register user"})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "user registered successfully"})
	}
}

func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := []models.User{}
		if err := db.Preload("Events").Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "could not get users",
			})
		}

		userResponses := make([]helpers.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = helpers.ToUserResponse(user)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "users fetched successfully",
			"data":    userResponses,
		})
	}
}

func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id cannot be empty"})
		}

		user := models.User{}
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "could not find user"})
		}

		if err := db.Delete(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not delete user"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user deleted successfully"})
	}
}

func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id cannot be empty"})
		}

		user := models.User{}
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "could not find user"})
		}

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		}

		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not update user"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user updated successfully"})
	}
}
