package middlewares

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/pkg/token"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func bNeedRefresh[T any](claim *token.RefreshTokenClaims[T]) bool {
	now := time.Now()

	expiry := claim.ExpiresAt.Time
	issuedAt := claim.IssuedAt.Time
	totalLifetime := expiry.Sub(issuedAt)
	remaining := expiry.Sub(now)
	remainingRatio := remaining.Seconds() / totalLifetime.Seconds()

	return remainingRatio <= token.REFRESH_TOKEN_ROTATION_THRESHOLD
}

func JwtAuth(cfg *config.Config) func(*gin.Context) {
	return jwtAuthImpl(cfg, true)
}

func OptionalJwtAuth(cfg *config.Config) func(*gin.Context) {
	return jwtAuthImpl(cfg, false)
}

func jwtAuthImpl(cfg *config.Config, mustAuth bool) func(*gin.Context) {
	return func(c *gin.Context) {
		if claim, ok := token.GetClaim[models.UserJwt](c, &token.AccessTokenClaims[models.UserJwt]{}, token.ACCESS_TOKEN_NAME, []byte(cfg.Opt.SecretKey)); ok {
			if uint(claim.Version) != uint(cfg.Opt.TokenVersion) {
				if mustAuth {
					ReturnUnauth(c, cfg)
					return
				}
				// Skip identification but continue for public access
				c.Next()
				return
			}

			// check if refresh token need to be rotated
			cookie, err := c.Cookie(token.REFRESH_TOKEN_NAME)
			if err == nil {
				refreshClaim, err := token.ValidateToken[models.UserJwtRefresh]([]byte(cfg.Opt.SecretKey), cookie, &token.RefreshTokenClaims[models.UserJwtRefresh]{})
				if err == nil && bNeedRefresh(refreshClaim) {
					var result struct {
						RefreshVersion int64 `json:"refresh_version"`
					}
					err = cfg.DB.Model(&models.User{}).
						Select(`refresh_version`).
						Where("id = ?", refreshClaim.Data.Id).
						First(&result).Error
					if err == nil && refreshClaim.RefreshVersion == uint(result.RefreshVersion) {
						refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.NewRefreshTokenClaim(refreshClaim.Data, refreshClaim.RefreshVersion, uint(cfg.Opt.TokenVersion)))
						refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Opt.SecretKey))
						if err == nil {
							SetCookie(c, token.REFRESH_TOKEN_NAME, refreshTokenString, token.REFRESH_TOKEN_LIFE_TIME, cfg.Opt.IsProd)
						}
					}
				}
			}

			c.Set(config.CONTEXT_USER, claim.Data)
			c.Next()
			return
		}

		cookie, err := c.Cookie(token.REFRESH_TOKEN_NAME)
		if err != nil {
			if mustAuth {
				ReturnUnauth(c, cfg)
				return
			}
			c.Next()
			return
		}

		refreshClaim, err := token.ValidateToken[models.UserJwtRefresh]([]byte(cfg.Opt.SecretKey), cookie, &token.RefreshTokenClaims[models.UserJwtRefresh]{})
		if err != nil {
			if mustAuth {
				ReturnUnauth(c, cfg)
				return
			}
			c.Next()
			return
		}

		if uint(refreshClaim.Version) != uint(cfg.Opt.TokenVersion) {
			if mustAuth {
				ReturnUnauth(c, cfg)
				return
			}
			c.Next()
			return
		}

		var result struct {
			RefreshVersion int64 `json:"refresh_version"`
		}
		err = cfg.DB.Model(&models.User{}).
			Select(`refresh_version`).
			Where("id = ?", refreshClaim.Data.Id).
			First(&result).Error
		if err != nil {
			if mustAuth {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					ReturnUnauth(c, cfg)
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Database error"})
				}
				return
			}
			c.Next()
			return
		}

		if refreshClaim.RefreshVersion != uint(result.RefreshVersion) {
			if mustAuth {
				ReturnUnauth(c, cfg)
				return
			}
			c.Next()
			return
		}

		var user struct {
			Id int64
		}
		if err := cfg.DB.Model(&models.User{}).Select("id").Where("id = ?", refreshClaim.Data.Id).First(&user).Error; err != nil {
			if mustAuth {
				ReturnUnauth(c, cfg)
				return
			}
			c.Next()
			return
		}

		accessTokenData := models.UserJwt{
			Id: user.Id,
		}
		accessTokenString, ok := token.RecreateAccessToken(accessTokenData, cfg.Opt.SecretKey, uint(cfg.Opt.TokenVersion))
		if ok {
			SetCookie(c, token.ACCESS_TOKEN_NAME, accessTokenString, token.ACCESS_TOKEN_LIFE_TIME, cfg.Opt.IsProd)
		}

		if bNeedRefresh(refreshClaim) {
			refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.NewRefreshTokenClaim(refreshClaim.Data, refreshClaim.RefreshVersion, uint(cfg.Opt.TokenVersion)))
			refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Opt.SecretKey))
			if err == nil {
				SetCookie(c, token.REFRESH_TOKEN_NAME, refreshTokenString, token.REFRESH_TOKEN_LIFE_TIME, cfg.Opt.IsProd)
			}
		}

		c.Set(config.CONTEXT_USER, accessTokenData)
		c.Next()
	}
}

func JwtSign(cfg *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		userData, ok := c.Get(config.CONTEXT_USER)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": "User data not found in context"})
			return
		}
		userJwt := userData.(models.UserJwt)

		var user models.User
		if err := cfg.DB.First(&user, userJwt.Id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": "User not found"})
			return
		}

		accessTokenString, ok := token.CreateAccessToken(userJwt, cfg.Opt.SecretKey, uint(cfg.Opt.TokenVersion))
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "error": "Failed to create access token"})
			return
		}
		SetCookie(c, token.ACCESS_TOKEN_NAME, accessTokenString, token.ACCESS_TOKEN_LIFE_TIME, cfg.Opt.IsProd)

		refreshVersionData, ok := c.Get(config.CONTEXT_REFRESH_VERSION)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": "refresh version not found"})
			return
		}
		refresh_version := refreshVersionData.(int64)

		refreshData := models.UserJwtRefresh{
			Id: userJwt.Id,
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token.NewRefreshTokenClaim(refreshData, uint(refresh_version), uint(cfg.Opt.TokenVersion)))
		refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Opt.SecretKey))
		if err != nil {
			log.Println("JwtSign error:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign jwt"})
			return
		}

		SetCookie(c, token.REFRESH_TOKEN_NAME, refreshTokenString, token.REFRESH_TOKEN_LIFE_TIME, cfg.Opt.IsProd)

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"user":   user,
			"token":  accessTokenString,
		})
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, ok := c.Get(config.CONTEXT_USER)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Unauthorized"})
			return
		}

		user := userData.(models.UserJwt)

		var dbUser models.User
		cfg := c.MustGet("config").(*config.Config)
		if err := cfg.DB.Select("is_admin").First(&dbUser, user.Id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Forbidden"})
			return
		}

		if !dbUser.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": false, "error": "Forbidden"})
			return
		}

		c.Next()
	}
}

func SetCookie(c *gin.Context, name string, value string, lifetime time.Duration, isProd bool) {
	c.SetCookie(name, value, int(lifetime.Seconds()), "/", "", isProd, true)
}

func ReturnUnauth(c *gin.Context, cfg *config.Config) {
	ClearToken(c, cfg)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "error": "Unauthorized"})
}

func GetUserFromContext(c *gin.Context) *models.UserJwt {
	userData, ok := c.Get(config.CONTEXT_USER)
	if !ok {
		return nil
	}
	user := userData.(models.UserJwt)
	return &user
}

func ClearToken(c *gin.Context, cfg *config.Config) {
	c.SetCookie(token.ACCESS_TOKEN_NAME, "", -1, "/", "", cfg.Opt.IsProd, true)
	c.SetCookie(token.REFRESH_TOKEN_NAME, "", -1, "/", "", cfg.Opt.IsProd, true)
}
