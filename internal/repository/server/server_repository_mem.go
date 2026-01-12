package repository

import (
	"errors"
	"rio/internal/models"
	"rio/internal/store"
)

type InMemoryServerRepository struct{}

func NewInMemoryServerRepository() *InMemoryServerRepository {
	return &InMemoryServerRepository{}
}

func (r *InMemoryServerRepository) Create(server *models.Server) error {
	for _, s := range store.Servers {
		if s.ULID == server.ULID {
			return errors.New("server with this ULID already exists")
		}
	}
	store.Servers = append(store.Servers, *server)
	return nil
}

func (r *InMemoryServerRepository) GetServerByID(ulid string) (*models.Server, error) {
	for i := range store.Servers {
		if store.Servers[i].ULID == ulid {
			return &store.Servers[i], nil
		}
	}
	return nil, errors.New("server not found")
}

func (r *InMemoryServerRepository) GetServersByUser(u_id string) ([]*models.Server, error) {
	var userServers []*models.Server

	for _, us := range store.UserServers {
		if us.UserID == u_id {
			for i := range store.Servers {
				if store.Servers[i].ULID == us.ServerID {
					userServers = append(userServers, &store.Servers[i])
					break
				}
			}
		}
	}

	return userServers, nil
}

func (r *InMemoryServerRepository) GetServerMembers(ulid string) ([]*models.User, error) {
	var members []*models.User

	for _, us := range store.UserServers {
		if us.ServerID == ulid {
			for i := range store.Users {
				if store.Users[i].ULID == us.UserID {
					members = append(members, &store.Users[i])
					break
				}
			}
		}
	}

	return members, nil
}

func (r *InMemoryServerRepository) UpdateServer(ulid string, server *models.Server) error {
	for i := range store.Servers {
		if store.Servers[i].ULID == ulid {
			store.Servers[i].Name = server.Name
			return nil
		}
	}
	return errors.New("server not found or no changes applied")
}

func (r *InMemoryServerRepository) DeleteServer(ulid string) error {
	for i := range store.Servers {
		if store.Servers[i].ULID == ulid {
			store.Servers = append(store.Servers[:i], store.Servers[i+1:]...)

			var newUserServers []models.UserServer
			for _, us := range store.UserServers {
				if us.ServerID != ulid {
					newUserServers = append(newUserServers, us)
				}
			}
			store.UserServers = newUserServers

			return nil
		}
	}
	return errors.New("server not found or already deleted")
}

func (r *InMemoryServerRepository) AddUserToServer(userID, serverID, role string) error {
	var userExists bool
	for _, u := range store.Users {
		if u.ULID == userID {
			userExists = true
			break
		}
	}
	if !userExists {
		return errors.New("user not found")
	}

	var serverExists bool
	for _, s := range store.Servers {
		if s.ULID == serverID {
			serverExists = true
			break
		}
	}
	if !serverExists {
		return errors.New("server not found")
	}

	for _, us := range store.UserServers {
		if us.UserID == userID && us.ServerID == serverID {
			return errors.New("user is already a member of this server")
		}
	}

	store.UserServers = append(store.UserServers, models.UserServer{
		UserID:   userID,
		ServerID: serverID,
		Role:     role,
	})

	return nil
}

func (r *InMemoryServerRepository) RemoveUserFromServer(userID, serverID string) error {
	for i, us := range store.UserServers {
		if us.UserID == userID && us.ServerID == serverID {
			store.UserServers = append(store.UserServers[:i], store.UserServers[i+1:]...)
			return nil
		}
	}
	return errors.New("membership not found (user may not be a member of the server)")
}

func (r *InMemoryServerRepository) UpdateUserRoleInServer(userID, serverID, newRole string) error {
	for i := range store.UserServers {
		if store.UserServers[i].UserID == userID && store.UserServers[i].ServerID == serverID {
			store.UserServers[i].Role = newRole
			return nil
		}
	}
	return errors.New("membership not found (user is not a member of this server, or user/server does not exist)")
}
