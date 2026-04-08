package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/pkg/password"
	"photoset/internal/repository"
	"time"
)

type UserService interface {
	Register(nickname, email, password string) (*domain.User, error)
	Login(email, password string) (*domain.User, error)
	GetProfile(userID uint) (*domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(nickname, email, password string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := password.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Nickname:     nickname,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         domain.RoleUser,
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !password.Check(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	return user, nil
}

func (s *userService) GetProfile(userID uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

