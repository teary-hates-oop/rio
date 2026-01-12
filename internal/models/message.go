package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	ULID      string `gorm:"type:varchar(26);primaryKey"`
	ChannelID string `gorm:"type:varchar(26);index"`
	UserID    string `gorm:"type:varchar(26);index"`
	Content   string `gorm:"not null"`
}
