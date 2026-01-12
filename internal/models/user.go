package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ULID     string   `gorm:"type:varchar(26);primaryKey;not null;unique"`
	Username string   `gorm:"size:255;not null;unique" json:"username"`
	Password string   `gorm:"size:255;not null" json:"-"`
	Role     string   `gorm:"size:20;default:'user'"`
	Servers  []Server `gorm:"many2many:user_servers;"`
}
