package models

import (
	"time"
)

type User struct {
	UserID       uint      `gorm:"primaryKey;autoIncrement;column:user_id"`
	Username     string    `gorm:"size:50;not null;column:username"`
	Email        string    `gorm:"size:100;unique;not null;column:email"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}
