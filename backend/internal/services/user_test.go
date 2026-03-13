package services

import (
	"backend/internal/mocks"
	"backend/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.Username == "testuser" && u.Name == "Test User" && u.IsAdmin == false
		})).Return(nil)

		err := service.CreateUser(ctx, "testuser", "password123", "Test User", false)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

		err := service.CreateUser(ctx, "testuser", "password123", "Test User", false)
		assert.Error(t, err)
		if err != nil {
			assert.Equal(t, "db error", err.Error())
		}
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Authenticate(t *testing.T) {
	ctx := context.Background()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

		authenticatedUser, err := service.Authenticate(ctx, "testuser", "password123")
		assert.NoError(t, err)
		assert.Equal(t, user.Username, authenticatedUser.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("GetByUsername", ctx, "testuser").Return(user, nil)

		authenticatedUser, err := service.Authenticate(ctx, "testuser", "wrongpassword")
		assert.Nil(t, authenticatedUser)
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, errors.New("not found"))

		authenticatedUser, err := service.Authenticate(ctx, "nonexistent", "password123")
		assert.Nil(t, authenticatedUser)
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdatePassword(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		Username:       "testuser",
		Password:       "oldhash",
		RefreshVersion: 0,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *models.User) bool {
			return u.ID == 1 && u.RefreshVersion == 1
		})).Return(nil)

		err := service.UpdatePassword(ctx, 1, "newpassword123")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewUserService(mockRepo)
		mockRepo.On("GetByID", ctx, int64(2)).Return(nil, errors.New("not found"))

		err := service.UpdatePassword(ctx, 2, "newpassword123")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
