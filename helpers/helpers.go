package helpers

import "github.com/ktariayman/go-api/models"

type UserResponse struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    EventIDs []uint `json:"events"`
}

type EventResponse struct {
    ID             uint   `json:"id"`
    Title          string `json:"title"`
    Description    string `json:"description"`
    Date           string `json:"date"`
    Location       string `json:"location"`
    UserID         uint   `json:"user_id"`
    ParticipantIDs []uint `json:"participants"`
}

func ToUserResponse(user models.User) UserResponse {
    eventIDs := make([]uint, len(user.Events))
    for i, event := range user.Events {
        eventIDs[i] = event.ID
    }
    return UserResponse{
        ID:       user.ID,
        Name:     user.Name,
        Email:    user.Email,
        EventIDs: eventIDs,
    }
}

func ToEventResponse(event models.Event) EventResponse {
    participantIDs := make([]uint, len(event.Participants))
    for i, user := range event.Participants {
        participantIDs[i] = user.ID
    }
    return EventResponse{
        ID:             event.ID,
        Title:          event.Title,
        Description:    event.Description,
        Date:           event.Date,
        Location:       event.Location,
        UserID:         event.UserID,
        ParticipantIDs: participantIDs,
    }
}
