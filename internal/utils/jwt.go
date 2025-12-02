package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/database"
	"github.com/riparuk/meet-book-api/internal/model"
	"github.com/riparuk/meet-book-api/internal/repository"
)

const (
	tokenExpiration = 24 * time.Hour
	tokenTypeBearer = "Bearer"
)

var (
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
	ErrInvalidClaims        = errors.New("invalid token claims")
	ErrInvalidUserID        = errors.New("invalid user ID in token")
	ErrUserNotFound         = errors.New("user not found")
)

// getJWTSecret retrieves the JWT secret from environment variables
func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}
	return []byte(secret), nil
}

// GenerateJWT creates a new JWT token for the given user ID and role
func GenerateJWT(userID string, role model.UserRole) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}

	if _, err := uuid.Parse(userID); err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	secret, err := getJWTSecret()
	if err != nil {
		return "", fmt.Errorf("failed to get JWT secret: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates the JWT token and returns the user ID and role if valid
func ValidateJWT(tokenString string) (string, model.UserRole, error) {
	if tokenString == "" {
		return "", "", ErrInvalidToken
	}

	// Remove 'Bearer ' prefix if present
	tokenString = strings.TrimPrefix(tokenString, tokenTypeBearer+" ")

	secret, err := getJWTSecret()
	if err != nil {
		return "", "", fmt.Errorf("failed to get JWT secret: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidSigningMethod, token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return "", "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", ErrInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", "", ErrInvalidUserID
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return "", "", fmt.Errorf("%w: invalid role type", ErrInvalidClaims)
	}

	return userID, model.UserRole(roleStr), nil
}

// GetUserAuth retrieves the authenticated user from the Gin context
func GetUserAuth(c *gin.Context) (*model.User, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists || userIDStr == "" {
		return nil, errors.New("user ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUserID, err)
	}

	userRepo := repository.NewUserRepository(database.DB)
	user, err := userRepo.FindByID(userID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
