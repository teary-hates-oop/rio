package repository

import (
	"errors"
	"rio/internal/db"
	"rio/internal/models"

	"github.com/jinzhu/gorm"
)

type DBUserRepository struct{}

func NewDBUserRepository() *DBUserRepository {
	return &DBUserRepository{}
}

func (r *DBUserRepository) Create(user *models.User) error {
	return db.DB.Create(user).Error
}

func (r *DBUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *DBUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := db.DB.Find(&users).Error
	return users, err
}

func (r *DBUserRepository) GetUserByID(id string) (*models.User, error) {
	var u models.User
	err := db.DB.Where("ul_id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	u.Password = ""

	return &u, nil
}
