package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
}

func MigrateEvents(db *gorm.DB) error {
	return db.AutoMigrate(&Event{})
}
