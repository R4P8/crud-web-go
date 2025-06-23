package config

import (
	"github.com/golang-jwt/jwt/v5"
)

var Jwt_Secret = []byte("SecretJWT29056!@#")

type ClaimsJWT struct {
	Username string `json:"Email"`
	jwt.RegisteredClaims
}

type ClaimsRefreshJWT struct {
	Username string `json:"Email"`
	jwt.RegisteredClaims
}
