package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN_NAME      = "access_token"
	ACCESS_TOKEN_LIFE_TIME = time.Minute * 15
)

type AccessTokenClaims[T interface{}] struct {
	TokenClaims[T]
	Version uint `json:"version"`
}

func NewAccessTokenClaim[T interface{}](data T, token_version uint) *AccessTokenClaims[T] {
	return &AccessTokenClaims[T]{
		TokenClaims: TokenClaims[T]{
			Data: data,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TOKEN_LIFE_TIME)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "property backend",
				Subject:   ACCESS_TOKEN_NAME,
			},
		},
		Version: token_version,
	}
}
