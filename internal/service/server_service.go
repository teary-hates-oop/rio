package service

import (
	"errors"
	"fmt"
	"html"
	"rio/internal/models"
	serverRepo "rio/internal/repository/server"
	userRepo "rio/internal/repository/user"
	"slices"
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

func validPermissions(userMembership, targetMembership models.UserServer) bool {
	if userMembership.UserID == targetMembership.UserID {
		return true
	}

	if userMembership.Role == "owner" {
		return true
	}

	switch userMembership.Role {
	case "admin":
		return targetMembership.Role == "moderator" || targetMembership.Role == "member"

	case "moderator":
		return targetMembership.Role == "member"

	case "member":
		return false

	default:
		return false
	}
}

func (s *ServerService) ListUserServers(currentUserID string) ([]*models.Server, error) {
	if currentUserID == "" {
		return nil, errors.New("current user ID is required")
	}

	servers, err := s.serverRepo.GetServersByUser(currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user servers: %w", err)
	}

	return servers, nil
}

func (s *ServerService) CreateServer(currentUserID, name string) (*models.Server, error) {

	name = html.EscapeString(strings.TrimSpace(name))

	if name == "" {
		return nil, errors.New("server name cannot be empty")
	}
	if len(name) < 3 || len(name) > 100 {
		return nil, errors.New("server name must be between 3 and 100 characters")
	}
	u, err := s.userRepo.GetUserByID(currentUserID)
	if err != nil {
		return nil, errors.New("user does not exist (contact dev)")
	}

	users := make([]models.User, 1)
	users = append(users, *u)
	newServer := models.Server{
		ULID:    ulid.Make().String(),
		Name:    name,
		OwnerID: currentUserID,
		Users:   users,
	}

	if newServer.ULID == "" {
		return nil, errors.New("failed to generate ULID")
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
	server, err := s.serverRepo.GetServerByID(serverID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.IsUserMember(currentUserID, serverID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a member of server")
	}
	return server, nil
}

func (s *ServerService) IsUserMember(currentUserID, serverID string) (bool, error) {
	membership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return false, err
	}
	return membership != nil, nil
}

func (s *ServerService) UpdateServerName(currentUserID, serverID, newName string) error {
	newName = html.EscapeString(strings.TrimSpace(newName))
	if newName == "" {
		return errors.New("server name cannot be empty")
	}
	if len(newName) < 3 || len(newName) > 100 {
		return errors.New("server name must be between 3 and 100 characters")
	}

	_, err := s.serverRepo.GetServerByID(serverID)
	if err != nil {
		return err
	}

	membership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return err
	}
	if membership == nil {
		return errors.New("you are not a member of this server")
	}

	if membership.Role != "owner" && membership.Role != "admin" {
		return errors.New("only server owner or admin can update the server name")
	}

	err = s.serverRepo.UpdateServer(serverID, &models.Server{Name: newName})
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerService) DeleteServer(currentUserID, serverID string) error {
	membership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return err
	}
	if membership == nil {
		return errors.New("user is not a member of server")
	}

	if membership.Role != "owner" {
		return errors.New("insufficient permissions")
	}

	err = s.serverRepo.DeleteServer(serverID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerService) AddMember(currentUserID, serverID, targetUserID, role string) error {
	callerMembership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return err
	}
	if callerMembership == nil {
		return errors.New("you are not a member of this server")
	}

	targetUser, err := s.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return errors.New("target user not found")
	}

	if currentUserID == targetUserID {
		return errors.New("cannot add yourself as a member")
	}

	if role == "owner" {
		return errors.New("there can only be one owner")
	}

	validRoles := []string{"admin", "moderator", "member"}
	if !slices.Contains(validRoles, role) {
		return errors.New("invalid role; must be one of: admin, moderator, member")
	}

	targetMembership, err := s.serverRepo.GetUserMembership(targetUserID, serverID)
	if err != nil {
		return err
	}

	if targetMembership != nil {
		if !validPermissions(*callerMembership, *targetMembership) {
			return errors.New("insufficient permissions to add or promote this user")
		}
	} else {
		if role == "admin" && callerMembership.Role != "owner" {
			return errors.New("only the server owner can invite admins")
		}
	}

	err = s.serverRepo.AddUserToServer(targetUserID, serverID, role)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerService) RemoveMember(currentUserID, serverID, targetUserID string) error {
	callerMembership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return err
	}
	if callerMembership == nil {
		return errors.New("you are not a member of this server")
	}

	targetUser, err := s.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return errors.New("target user not found")
	}

	targetMembership, err := s.serverRepo.GetUserMembership(targetUserID, serverID)
	if err != nil {
		return err
	}

	if targetMembership == nil {
		return errors.New("target user is not a member of this server")
	}

	if targetMembership.Role == "owner" {
		return errors.New("the server owner cannot be removed")
	}

	if currentUserID == targetUserID && callerMembership.Role == "owner" {
		return errors.New("the server owner cannot remove themselves")
	}

	if !validPermissions(*callerMembership, *targetMembership) {
		return errors.New("insufficient permissions to remove this member")
	}

	err = s.serverRepo.RemoveUserFromServer(targetUserID, serverID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServerService) ChangeMemberRole(currentUserID, serverID, targetUserID, role string) error {
	callerMembership, err := s.serverRepo.GetUserMembership(currentUserID, serverID)
	if err != nil {
		return err
	}
	if callerMembership == nil {
		return errors.New("you are not a member of this server")
	}

	targetMembership, err := s.serverRepo.GetUserMembership(targetUserID, serverID)
	if err != nil {
		return err
	}
	if targetMembership == nil {
		return errors.New("target user is not a member of this server")
	}

	if currentUserID == targetUserID && callerMembership.Role == "owner" {
		return errors.New("the server owner cannot change their own role through this endpoint")
	}

	if !validPermissions(*callerMembership, *targetMembership) {
		return errors.New("insufficient permissions to change this user's role")
	}

	if role == "owner" {
		return errors.New("ownership cannot be delegated through role changes; use a dedicated transfer endpoint")
	}

	validRoles := []string{"admin", "moderator", "member"}
	if !slices.Contains(validRoles, role) {
		return errors.New("invalid role; must be one of: admin, moderator, member")
	}

	err = s.serverRepo.UpdateUserRoleInServer(targetUserID, serverID, role)
	if err != nil {
		return err
	}

	return nil
}
