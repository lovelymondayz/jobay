package utils

import (
	"os"
	"time"
	"wifi-tracker-be/src/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(user *models.Users) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.Name,
        "role_name": user.Role.Name,
        "client_name": user.Client.Name,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseJWT(tokenStr string) (uint, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return uint(claims["user_id"].(float64)), nil
    }

    return 0, err
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
