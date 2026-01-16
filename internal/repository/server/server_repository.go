package repository

import "rio/internal/models"

type ServerRepository interface {
	Create(server *models.Server) error
	CreateMembership(membership *models.UserServer) error
	GetUserMembership(u_id, s_id string) (*models.UserServer, error)
	GetServerByID(ulid string) (*models.Server, error)
	GetServersByUser(u_id string) ([]*models.Server, error)
	GetServerMembers(ulid string) ([]*models.User, error)
	UpdateServer(ulid string, server *models.Server) error
	DeleteServer(ulid string) error
	AddUserToServer(userID, serverID, role string) error
	RemoveUserFromServer(userID, serverID string) error
	UpdateUserRoleInServer(userID, serverID, newRole string) error
}
