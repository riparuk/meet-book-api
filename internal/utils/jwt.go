package utils

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/database"
	"github.com/riparuk/meet-book-api/internal/model"
	"github.com/riparuk/meet-book-api/internal/repository"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id in token")
	}

	return userID, nil
}

func GetUserAuth(c *gin.Context) (*model.User, error) {
	userRepo := repository.NewUserRepository(database.DB)

	userIDStr, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("unauthorized")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := userRepo.FindByID(userID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
