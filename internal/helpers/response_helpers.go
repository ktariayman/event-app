package helpers

import "github.com/ktariayman/go-api/internal/models"

type UserResponse struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    EventIDs []uint `json:"events"`
    TotalVotes  int              `json:"total_votes"`
    EventVotes  []EventVoteCount `json:"event_votes"`
}

type EventVoteCount struct {
    EventID uint `json:"event_id"`
    Votes   int  `json:"votes"`
}

type EventResponse struct {
    ID             uint   `json:"id"`
    Title          string `json:"title"`
    Description    string `json:"description"`
    Date           string `json:"date"`
    Location       string `json:"location"`
    UserID         uint   `json:"user_id"`
    ParticipantIDs []uint `json:"participants"`
    Votes          int    `json:"votes"`  
}

func ToUserResponse(user models.User) UserResponse {
    eventIDs := make([]uint, len(user.Events))
    eventVotes := make([]EventVoteCount, len(user.Events))
    totalVotes := 0

    for i, event := range user.Events {
        eventIDs[i] = event.ID
        eventVotes[i] = EventVoteCount{
            EventID: event.ID,
            Votes:   event.Votes,
        }
        totalVotes += event.Votes
    }

    return UserResponse{
        ID:         user.ID,
        Name:       user.Name,
        Email:      user.Email,
        EventIDs:   eventIDs,
        TotalVotes: totalVotes,
        EventVotes: eventVotes,
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
        Votes:          event.Votes,
    }
}
