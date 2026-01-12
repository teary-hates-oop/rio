package repository

import "rio/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindAll() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
}
