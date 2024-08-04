package models

import "gorm.io/gorm"

type Event struct {
    gorm.Model
    Title        string `json:"title" gorm:"not null"`
    Description  string `json:"description" gorm:"not null"`
    Date         string `json:"date" gorm:"not null"`
    Location     string `json:"location" gorm:"not null"`
    UserID       uint   `json:"user_id" gorm:"not null"`
    Participants []User `gorm:"many2many:event_participants;" json:"participants"`
    Votes        int    `json:"votes" gorm:"default:0"`
}

func (event *Event) BeforeDelete(tx *gorm.DB) (err error) {
    err = tx.Model(event).Association("Participants").Clear()
    return
}

type EventVote struct {
    gorm.Model
    EventID uint `json:"event_id" gorm:"not null;index"`
    UserID  uint `json:"user_id" gorm:"not null;index"`
    Vote    int  `json:"vote"` 
}

type VoteRequest struct {
    Action int `json:"action"`
}

const (
    VoteDown = 0
    VoteUp   = 1
)

func MigrateEvents(db *gorm.DB) error {
    err := db.AutoMigrate(&Event{})
    if err != nil {
        return err
    }
    return db.AutoMigrate(&EventVote{})
}
