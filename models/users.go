package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`
    Events   []Event `gorm:"many2many:event_participants;" json:"events"`
}

func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
    err = tx.Model(user).Association("Events").Clear()
    return
}

func MigrateUsers(db *gorm.DB) error {
    return db.AutoMigrate(&User{})
}
