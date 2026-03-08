package graph

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("no gin context", func(t *testing.T) {
		user := GetUser(context.Background())
		assert.Nil(t, user)
	})

	t.Run("no user in context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx := context.WithValue(context.Background(), ginContextKey, c)
		user := GetUser(ctx)
		assert.Nil(t, user)
	})

	t.Run("user as value in context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedUser := models.UserJwt{Id: 123}
		c.Set(config.CONTEXT_USER, expectedUser)
		ctx := context.WithValue(context.Background(), ginContextKey, c)
		
		user := GetUser(ctx)
		assert.NotNil(t, user)
		assert.Equal(t, int64(123), user.Id)
	})

	t.Run("user as pointer in context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedUser := &models.UserJwt{Id: 456}
		c.Set(config.CONTEXT_USER, expectedUser)
		ctx := context.WithValue(context.Background(), ginContextKey, c)
		
		user := GetUser(ctx)
		assert.NotNil(t, user)
		assert.Equal(t, int64(456), user.Id)
	})
}

func TestIsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("no user", func(t *testing.T) {
		isAdmin := IsAdmin(context.Background(), nil)
		assert.False(t, isAdmin)
	})

	t.Run("no config", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(config.CONTEXT_USER, models.UserJwt{Id: 1})
		ctx := context.WithValue(context.Background(), ginContextKey, c)
		
		isAdmin := IsAdmin(ctx, nil)
		assert.False(t, isAdmin)
	})
}
