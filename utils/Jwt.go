package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type JWTClaim struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(email, role string) (string, error) {

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	now := time.Now()
	claims := &JWTClaim{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
