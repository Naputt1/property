package token

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims[T interface{}] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

type ExpireError struct {
	Message string
}

func (e *ExpireError) Error() string {
	return e.Message
}

type IClaim[T interface{}] interface {
	jwt.Claims
	*AccessTokenClaims[T] | *RefreshTokenClaims[T]
}

func ValidateToken[T interface{}, C IClaim[T]](secretKey []byte, cookie string, claims C) (C, error) {
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(C)
	if !ok || !token.Valid {
		return nil, err
	}

	expiteAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if expiteAt.Before(time.Now()) {
		return nil, &ExpireError{}
	}

	return claims, nil
}

func RecreateAccessToken[T interface{}](data T, secretKey string, version uint) (string, bool) {
	return CreateAccessToken(data, secretKey, version)
}

func CreateAccessToken[T interface{}](data T, secretKey string, version uint) (string, bool) {
	claim := NewAccessTokenClaim(data, version)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", false
	}

	return accessTokenString, true
}

func GetClaim[T interface{}, C IClaim[T]](c *gin.Context, claims C, cookieName string, secretKey []byte) (C, bool) {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]
			if tokenString != "null" && tokenString != "" {
				claim, err := ValidateToken[T](secretKey, tokenString, claims)
				if err == nil {
					return claim, true
				}
				log.Printf("get claim validate header: %s", err.Error())
			}
		}
	}

	// Fallback to cookie
	cookie, err := c.Cookie(cookieName)
	if err == nil {
		claim, err := ValidateToken[T](secretKey, cookie, claims)
		if err == nil {
			return claim, true
		}

		log.Printf("get claim validate cookie: %s - %s", cookieName, err.Error())
		return nil, false
	}

	return nil, false
}
