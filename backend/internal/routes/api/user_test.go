package api

import (
	"backend/internal/mocks"
	"backend/internal/models"
	"backend/internal/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUserSvc := new(mocks.MockUserService)
		svcs := &services.Services{User: mockUserSvc}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := &models.User{ID: 1, Username: "testuser"}
		c.Set("services", svcs)
		c.Set("user", models.UserJwt{Id: 1})

		reqBody := ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/user/change-password", bytes.NewBuffer(jsonBody))

		mockUserSvc.On("GetUserByID", mock.Anything, int64(1)).Return(user, nil)
		mockUserSvc.On("Authenticate", mock.Anything, "testuser", "oldpassword").Return(user, nil)
		mockUserSvc.On("UpdatePassword", mock.Anything, int64(1), "newpassword123").Return(nil)

		ChangePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response BaseResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Status)
		assert.Equal(t, "password updated successfully", response.Message)

		mockUserSvc.AssertExpectations(t)
	})

	t.Run("invalid old password", func(t *testing.T) {
		mockUserSvc := new(mocks.MockUserService)
		svcs := &services.Services{User: mockUserSvc}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := &models.User{ID: 1, Username: "testuser"}
		c.Set("services", svcs)
		c.Set("user", models.UserJwt{Id: 1})

		reqBody := ChangePasswordRequest{
			OldPassword: "wrongpassword",
			NewPassword: "newpassword123",
		}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/user/change-password", bytes.NewBuffer(jsonBody))

		mockUserSvc.On("GetUserByID", mock.Anything, int64(1)).Return(user, nil)
		mockUserSvc.On("Authenticate", mock.Anything, "testuser", "wrongpassword").Return(nil, assert.AnError)

		ChangePassword(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUserSvc.AssertExpectations(t)
	})
}
