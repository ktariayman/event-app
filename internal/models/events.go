package models

import "gorm.io/gorm"

type Event struct {
    gorm.Model
    Title        string `json:"title"`
    Description  string `json:"description"`
    Date         string `json:"date"`
    Location     string `json:"location"`
    UserID       uint   `json:"user_id"`
    Participants []User `gorm:"many2many:event_participants;" json:"participants"`
}

func (event *Event) BeforeDelete(tx *gorm.DB) (err error) {
    err = tx.Model(event).Association("Participants").Clear()
    return
}

func MigrateEvents(db *gorm.DB) error {
    return db.AutoMigrate(&Event{})
}
