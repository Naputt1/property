package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	REFRESH_TOKEN_NAME               = "refresh_token"
	REFRESH_TOKEN_LIFE_TIME          = time.Hour * 24 * 3
	REFRESH_TOKEN_ROTATION_THRESHOLD = 0.2
)

type RefreshTokenClaims[T interface{}] struct {
	TokenClaims[T]
	RefreshVersion uint `json:"refresh_version"`
	Version        uint `json:"version"`
}

func NewRefreshTokenClaim[T interface{}](data T, refresh_token uint, token_version uint) *RefreshTokenClaims[T] {
	return &RefreshTokenClaims[T]{
		TokenClaims: TokenClaims[T]{
			Data: data,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(REFRESH_TOKEN_LIFE_TIME)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "property backend",
				Subject:   REFRESH_TOKEN_NAME,
			},
		},
		RefreshVersion: refresh_token,
		Version:        token_version,
	}
}
