package repository

import (
	"errors"
	"rio/internal/db"
	"rio/internal/models"

	"github.com/jinzhu/gorm"
)

type DBServerRepository struct{}

func NewDBServerRepository() *DBServerRepository {
	return &DBServerRepository{}
}

func (r *DBServerRepository) Create(server *models.Server) error {
	if server.ULID == "" {
		return errors.New("server ULID is empty")
	}
	return db.DB.Create(server).Error
}

func (r *DBServerRepository) CreateMembership(membership *models.UserServer) error {
	return db.DB.Create(membership).Error
}

func (r *DBServerRepository) GetUserMembership(u_id, s_id string) (*models.UserServer, error) {
	var membership models.UserServer
	err := db.DB.
		Where("user_id = ? AND server_id = ?", u_id, s_id).
		First(&membership).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &membership, nil
}

func (r *DBServerRepository) GetServerByID(ulid string) (*models.Server, error) {
	var s models.Server
	err := db.DB.Where("ul_id = ?", ulid).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil

}

func (r *DBServerRepository) GetServersByUser(u_id string) ([]*models.Server, error) {
	var servers []*models.Server

	err := db.DB.
		Joins("JOIN user_servers ON user_servers.server_id = servers.ul_id").
		Where("user_servers.user_id = ?", u_id).
		Find(&servers).Error

	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (r *DBServerRepository) GetServerMembers(ulid string) ([]*models.User, error) {
	var members []*models.User

	err := db.DB.
		Joins("JOIN user_servers ON user_servers.user_id = users.ul_id").
		Where("user_servers.server_id = ?", ulid).
		Find(&members).Error

	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *DBServerRepository) UpdateServer(ulid string, server *models.Server) error {
	result := db.DB.Model(&models.Server{}).
		Where("ul_id = ?", ulid).
		Update("name", server.Name)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("server not found or no changes applied")
	}

	return nil
}

func (r *DBServerRepository) DeleteServer(ulid string) error {
	result := db.DB.Where("ul_id = ?", ulid).Delete(&models.Server{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("server not found or already deleted")
	}

	return nil
}

func (r *DBServerRepository) AddUserToServer(userID, serverID, role string) error {
	var userCount, serverCount int64
	db.DB.Model(&models.User{}).Where("ul_id = ?", userID).Count(&userCount)
	db.DB.Model(&models.Server{}).Where("ul_id = ?", serverID).Count(&serverCount)

	if userCount == 0 {
		return errors.New("user not found")
	}
	if serverCount == 0 {
		return errors.New("server not found")
	}

	var existingCount int64
	db.DB.Model(&models.UserServer{}).
		Where("user_id = ? AND server_id = ?", userID, serverID).
		Count(&existingCount)

	if existingCount > 0 {
		return errors.New("user is already a member of this server")
	}

	membership := models.UserServer{
		UserID:   userID,
		ServerID: serverID,
		Role:     role,
	}

	return db.DB.Create(&membership).Error
}

func (r *DBServerRepository) RemoveUserFromServer(userID, serverID string) error {
	result := db.DB.
		Where("user_id = ? AND server_id = ?", userID, serverID).
		Delete(&models.UserServer{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("membership not found (user may not be a member of the server)")
	}

	return nil
}

func (r *DBServerRepository) UpdateUserRoleInServer(userID, serverID, newRole string) error {
	result := db.DB.Model(&models.UserServer{}).
		Where("user_id = ? AND server_id = ?", userID, serverID).
		Update("role", newRole)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("membership not found (user is not a member of this server, or user/server does not exist)")
	}

	return nil
}
