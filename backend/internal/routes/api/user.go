package api

import (
	"backend/internal/config"
	"backend/internal/routes/middlewares"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword godoc
// @Summary Change own password
// @Description Change the password of the currently logged-in user.
// @Tags user
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param body body ChangePasswordRequest true "Password details"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse
// @Router /user/change-password [post]
func ChangePassword(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)
	userClaims := middlewares.GetUserFromContext(c)
	if userClaims == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Status: false, Error: "unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	user, err := svcs.User.GetUserByID(c.Request.Context(), userClaims.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "user not found"})
		return
	}

	// Authenticate with old password
	_, err = svcs.User.Authenticate(c.Request.Context(), user.Username, req.OldPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Status: false, Error: "invalid old password"})
		return
	}

	// Update password
	err = svcs.User.UpdatePassword(c.Request.Context(), user.ID, req.NewPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{Status: true, Message: "password updated successfully"})
}

func RegisterUserRoutes(r *gin.RouterGroup, cfg *config.Config, svcs *services.Services) {
	r.Use(func(c *gin.Context) {
		c.Set("services", svcs)
		c.Next()
	})
	r.POST("/change-password", ChangePassword)
}
