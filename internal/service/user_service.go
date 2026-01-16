package service

import (
	"errors"
	"html"
	"strings"

	"rio/internal/models"
	repository "rio/internal/repository/user"
	"rio/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func validateUsername(username *string) error {
	clean := html.EscapeString(strings.TrimSpace(*username))
	if clean == "" {
		return errors.New("username cannot be empty")
	}

	if len(clean) < 3 || len(clean) > 25 {
		return errors.New("username must be between 3 and 25 characters")
	}

	first := clean[0]
	last := clean[len(clean)-1]

	if first == '_' || last == '_' {
		return errors.New("username cannot start or end with an underscore")
	}

	for _, ch := range clean {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-') {
			return errors.New("username can only contain letters, numbers, underscores, and hyphens")
		}
	}

	*username = clean

	return nil
}

func (s *UserService) LoginCheck(username, password string) (string, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("invalid username or password")
	}

	if err := VerifyPassword(password, u.Password); err != nil {
		return "", errors.New("invalid username or password")
	}

	return token.GenerateToken(u.ULID)
}

func (s *UserService) Register(username, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	if err := validateUsername(&username); err != nil {
		return nil, err
	}

	user := &models.User{
		ULID:     ulid.Make().String(),
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindCurrentUser(c *gin.Context) (models.User, error) {
	uid, err := token.ExtractTokenID(c)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.repo.GetUserByID(uid)
	if err != nil {
		return models.User{}, err
	}
	user.Password = ""
	return *user, nil
}

func (s *UserService) FindUser(username string) (*models.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}
