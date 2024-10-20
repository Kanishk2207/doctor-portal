package utils

import (
	"auth_service/internal/configs"
	"auth_service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(user *models.User) (string, error) {
	config := configs.LoadConfig()

	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     int(GetCurrentUnixTime()) + config.JWTExpiry,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}
