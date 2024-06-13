package models

import (
	"time"
)

type Group struct {
	GroupID   uint      `gorm:"primaryKey;autoincrement;column:group_id"`
	GroupName string    `gorm:"size:100;not null;column:group_name"`
	AdminID   uint      `gorm:"not null;column:admin_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
}

type Groupmember struct {
	GroupID  uint      `grom:"primaryKey;column:group_id"`
	UserID   uint      `grom:"primaryKey;column:user_id"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;column:joined_at"`
}
