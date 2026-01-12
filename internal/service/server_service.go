package service

import (
	"errors"
	"fmt"
	"html"
	"rio/internal/models"
	serverRepo "rio/internal/repository/server"
	userRepo "rio/internal/repository/user"
	"strings"

	"github.com/oklog/ulid/v2"
)

type ServerService struct {
	serverRepo serverRepo.ServerRepository
	userRepo   userRepo.UserRepository
}

func NewServerService(
	sRepo serverRepo.ServerRepository,
	uRepo userRepo.UserRepository,
) *ServerService {
	return &ServerService{
		serverRepo: sRepo,
		userRepo:   uRepo,
	}
}

func (s *ServerService) CreateServer(currentUserID, name string) (*models.Server, error) {

	name = html.EscapeString(strings.TrimSpace(name))

	if name == "" {
		return nil, errors.New("server name cannot be empty")
	}
	if len(name) < 3 || len(name) > 100 {
		return nil, errors.New("server name must be between 3 and 100 characters")
	}

	newServer := models.Server{
		ULID:    ulid.Make().String(),
		Name:    name,
		OwnerID: currentUserID,
	}

	if err := s.serverRepo.Create(&newServer); err != nil {
		return nil, fmt.Errorf("failed to create server: %w, called by user: %v", err, currentUserID)
	}

	membership := models.UserServer{
		UserID:   currentUserID,
		ServerID: newServer.ULID,
		Role:     "owner",
	}

	if err := s.serverRepo.CreateMembership(&membership); err != nil {
		return nil, fmt.Errorf("failed to assign owner role: %w for sID: %v, uID: %v", err, newServer.ULID, currentUserID)
	}

	return &newServer, nil
}

func (s *ServerService) GetServer(currentUserID, serverID string) (*models.Server, error) {

}
