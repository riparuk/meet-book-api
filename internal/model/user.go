package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"` // don't expose password in JSON
	Role      UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Name     string   `json:"name" binding:"required" example:"Rifa Faruqi"`
	Email    string   `json:"email" binding:"required,email" example:"riparuk@gmail.com"`
	Password string   `json:"password" binding:"required" example:"strongpassword"`
	Role     UserRole `json:"role" binding:"required" example:"user"`
}

type UpdateUserInput struct {
	Name     string   `json:"name" binding:"required" example:"Rifa Faruqi"`
	Email    string   `json:"email" binding:"required,email" example:"riparuk@gmail.com"`
	Password string   `json:"password" binding:"required" example:"strongpassword"`
	Role     UserRole `json:"role" binding:"required" example:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"riparuk@gmail.com"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"Rifa Faruqi"`
	Email    string `json:"email" binding:"required,email" example:"riparuk@gmail.com"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}
