package database

import (
	"fmt"
	"log"

	"github.com/ktariayman/go-api/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func NewConnection(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return nil, err
	}
	return db, nil
}