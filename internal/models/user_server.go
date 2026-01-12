package models

import (
	"time"
)

type UserServer struct {
	UserID   string    `gorm:"primaryKey;type:varchar(26)"`
	ServerID string    `gorm:"primaryKey;type:varchar(26)"`
	Role     string    `gorm:"type:varchar(20);not null;default:'member'"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
}
