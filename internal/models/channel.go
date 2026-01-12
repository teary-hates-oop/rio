package models

import (
	"github.com/jinzhu/gorm"
)

type Channel struct {
	gorm.Model
	ULID     string `gorm:"type:varchar(26);primaryKey"`
	ServerID string `gorm:"type:varchar(26);index"`
	Name     string `gorm:"not null"`
}
