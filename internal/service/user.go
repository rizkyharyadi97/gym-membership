package service

import (
	"context"
	"errors"
	"strings"

	"gym-membership/internal/model"
	"gym-membership/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRequired   = errors.New("username is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrUserExists         = errors.New("username or email already exists")
	ErrInvalidDeposit     = errors.New("deposit amount must be greater than zero")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	AddDeposit(ctx context.Context, userID uint, amount float64) error
	Login(ctx context.Context, email string, password string) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Register(ctx context.Context, user *model.User) error {

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	if user.Username == "" {
		return ErrUsernameRequired
	}

	if user.Email == "" {
		return ErrEmailRequired
	}

	if user.Password == "" {
		return ErrPasswordRequired
	}

	// Check duplicate username.
	existingUser, err := s.userRepository.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return ErrUserExists
	}

	if err != nil &&
		!errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	// Check duplicate email.
	existingUser, err = s.userRepository.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return ErrUserExists
	}

	if err != nil &&
		!errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	// Hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.DepositAmount = 0

	return s.userRepository.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.userRepository.GetAll(ctx)
}

func (s *userService) AddDeposit(ctx context.Context, userID uint, amount float64) error {
	if amount <= 0 {
		return ErrInvalidDeposit
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.userRepository.UpdateDeposit(ctx, userID, amount)
}

func (s *userService) Login(ctx context.Context, email string, password string) (*model.User, error) {

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
