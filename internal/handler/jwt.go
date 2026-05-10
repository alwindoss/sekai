package handler

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret() ([]byte, error) {
	jwtSecret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		return nil, fmt.Errorf("Env variable JWT_SECRET not found")
	}
	return []byte(jwtSecret), nil
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
