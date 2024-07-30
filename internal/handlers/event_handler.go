package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/ktariayman/go-api/internal/helpers"
	"github.com/ktariayman/go-api/internal/models"
	"gorm.io/gorm"
)

func CreateEvent(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(float64)
		event := models.Event{
			UserID: uint(userID),
		}
		if err := c.BodyParser(&event); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		}
		if err := db.Create(&event).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "could not create event"})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "event has been added"})
	}
}

func DeleteEvent(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(float64)
		id := c.Params("id")
		if id == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "id cannot be empty"})
		}

		event := models.Event{}
		if err := db.Where("id = ? AND user_id = ?", id, uint(userID)).First(&event).Error; err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "could not find event"})
		}

		if err := db.Delete(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "could not delete event"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "event deleted successfully"})
	}
}

func GetEvents(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		events := []models.Event{}
		if err := db.Preload("Participants").Find(&events).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "could not get events",
			})
		}

		eventResponses := make([]helpers.EventResponse, len(events))
		for i, event := range events {
			eventResponses[i] = helpers.ToEventResponse(event)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "events fetched successfully",
			"data":    eventResponses,
		})
	}
}

func GetEventByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "id cannot be empty"})
		}

		event := models.Event{}
		if err := db.Where("id = ?", id).First(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "could not get the event"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "event fetched successfully", "data": event})
	}
}

func UpdateEvent(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(float64)
		id := c.Params("id")
		event := models.Event{}
		if err := db.Where("id = ? AND user_id = ?", id, uint(userID)).First(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "event not found"})
		}
		if err := c.BodyParser(&event); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		}
		if err := db.Save(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "could not update event"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "event updated successfully"})
	}
}

func ParticipateInEvent(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := uint(c.Locals("userID").(float64))
		eventID := c.Params("id")
		if eventID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Event ID cannot be empty"})
		}

		event := models.Event{}
		if err := db.Preload("Participants").Where("id = ?", eventID).First(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Event not found"})
		}

		user := models.User{}
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not found"})
		}

		// Check if the user is already a participant
		isParticipant := false
		for _, participant := range event.Participants {
			if participant.ID == userID {
				isParticipant = true
				break
			}
		}

		if isParticipant {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "You are already a participant in this event"})
		}

		if err := db.Model(&event).Association("Participants").Append(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not participate in event"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully participated in event"})
	}
}

func CancelParticipation(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := uint(c.Locals("userID").(float64))
		eventID := c.Params("id")
		if eventID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Event ID cannot be empty"})
		}

		event := models.Event{}
		if err := db.Preload("Participants").Where("id = ?", eventID).First(&event).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Event not found"})
		}

		user := models.User{}
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User not found"})
		}

		// Check if the user is a participant
		isParticipant := false
		for _, participant := range event.Participants {
			if participant.ID == userID {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "You are not a participant in this event"})
		}

		if err := db.Model(&event).Association("Participants").Delete(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not cancel participation in event"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Successfully canceled participation in event"})
	}
}
