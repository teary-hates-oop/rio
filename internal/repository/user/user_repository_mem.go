package repository

import (
	"errors"
	"rio/internal/models"
	"rio/internal/store"
	"strings"
)

type InMemoryUserRepository struct{}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{}
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
	for _, u := range store.Users {
		if strings.EqualFold(u.Username, user.Username) {
			return errors.New("username already taken")
		}
		if u.ULID == user.ULID {
			return errors.New("user ID already exists")
		}
	}
	store.Users = append(store.Users, *user)
	return nil
}

func (r *InMemoryUserRepository) FindByUsername(username string) (*models.User, error) {
	for _, u := range store.Users {
		if strings.EqualFold(u.Username, username) {
			return &u, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) FindAll() ([]models.User, error) {
	return store.Users, nil
}
