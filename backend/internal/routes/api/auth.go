package api

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/routes/middlewares"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary User login
// @Description Log in a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginBody true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)
	var body LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	username := strings.ToLower(body.Username)

	var user models.User
	if err := cfg.DB.Select("id, password, refresh_version").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Status: false, Error: "invalid username or password"})
			return
		}

		log.Println("failed login: ", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "Database error"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Status: false, Error: "invalid username or password"})
		return
	}

	c.Set(config.CONTEXT_USER, models.UserJwt{
		Id: user.ID,
	})
	c.Set(config.CONTEXT_REFRESH_VERSION, user.RefreshVersion)

	c.Next()
}

// Logout godoc
// @Summary User logout
// @Description Log out the current user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} BaseResponse
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)
	middlewares.ClearToken(c, cfg)

	c.JSON(http.StatusOK, BaseResponse{
		Status: true,
	})
}

func RegisterAuthRoutes(r *gin.RouterGroup, cfg *config.Config) {
	r.POST("/login", Login, middlewares.JwtSign(cfg))
	r.POST("/logout", Logout)
}
