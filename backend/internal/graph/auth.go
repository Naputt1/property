package graph

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const ginContextKey contextKey = "ginContext"

// GinContextToContextMiddleware puts the gin.Context into the request context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetGinContext extracts the gin.Context from the request context
func GetGinContext(ctx context.Context) *gin.Context {
	ginCtx := ctx.Value(ginContextKey)
	if ginCtx == nil {
		return nil
	}
	return ginCtx.(*gin.Context)
}

// GetUser extracts the user from the gin.Context inside the request context
func GetUser(ctx context.Context) *models.UserJwt {
	c := GetGinContext(ctx)
	if c == nil {
		return nil
	}
	user, ok := c.Get(config.CONTEXT_USER)
	if !ok {
		return nil
	}

	if u, ok := user.(models.UserJwt); ok {
		return &u
	}
	if u, ok := user.(*models.UserJwt); ok {
		return u
	}

	return nil
}

// IsAdmin checks if the current user is an admin
func IsAdmin(ctx context.Context, cfg *config.Config) bool {
	user := GetUser(ctx)
	if user == nil {
		return false
	}

	if cfg == nil || cfg.DB == nil {
		return false
	}

	var dbUser models.User
	if err := cfg.DB.Model(&models.User{}).Select("is_admin").Where("id = ?", user.Id).First(&dbUser).Error; err != nil {
		return false
	}

	return dbUser.IsAdmin
}
