package models

import (
	"github.com/jinzhu/gorm"
)

type Server struct {
	gorm.Model
	ULID     string `gorm:"type:varchar(26);primaryKey;unique"`
	Name     string `gorm:"size:255;not null"`
	OwnerID  string `gorm:"type:varchar(26);index"`
	Users    []User `gorm:"many2many:user_servers;"`
	Channels []Channel
}
