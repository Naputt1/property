package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, username, password, name string, isAdmin bool) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Name:     name,
		IsAdmin:  isAdmin,
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, username, name string, isAdmin bool) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Username = username
	user.Name = name
	user.IsAdmin = isAdmin

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	// Increment refresh version to invalidate existing tokens
	user.RefreshVersion++

	return s.repo.Update(ctx, user)
}
