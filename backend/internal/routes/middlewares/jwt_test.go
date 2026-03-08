package middlewares

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/pkg/token"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOptionalJwtAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "testsecret"
	cfg := &config.Config{
		Opt: config.Option{
			SecretKey:    secret,
			TokenVersion: 1,
		},
	}

	t.Run("no token allows request", func(t *testing.T) {
		r := gin.New()
		r.Use(OptionalJwtAuth(cfg))
		r.GET("/test", func(c *gin.Context) {
			_, exists := c.Get(config.CONTEXT_USER)
			assert.False(t, exists)
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("valid token in header identifies user", func(t *testing.T) {
		user := models.UserJwt{Id: 1}
		tokenString, _ := token.CreateAccessToken(user, secret, 1)

		r := gin.New()
		r.Use(OptionalJwtAuth(cfg))
		r.GET("/test", func(c *gin.Context) {
			val, exists := c.Get(config.CONTEXT_USER)
			assert.True(t, exists)
			assert.Equal(t, user.Id, val.(models.UserJwt).Id)
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid token allows request but no user", func(t *testing.T) {
		r := gin.New()
		r.Use(OptionalJwtAuth(cfg))
		r.GET("/test", func(c *gin.Context) {
			_, exists := c.Get(config.CONTEXT_USER)
			assert.False(t, exists)
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestJwtAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "testsecret"
	cfg := &config.Config{
		Opt: config.Option{
			SecretKey:    secret,
			TokenVersion: 1,
		},
	}

	t.Run("no token aborts request", func(t *testing.T) {
		r := gin.New()
		r.Use(JwtAuth(cfg))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("valid token in header allows request", func(t *testing.T) {
		user := models.UserJwt{Id: 1}
		tokenString, _ := token.CreateAccessToken(user, secret, 1)

		r := gin.New()
		r.Use(JwtAuth(cfg))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
