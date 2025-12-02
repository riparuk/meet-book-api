package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string         `json:"name"`
	Capacity  int            `json:"capacity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateRoomInput struct {
	Name     string `json:"name" binding:"required" example:"Meeting Room 1"`
	Capacity int    `json:"capacity" binding:"required" example:"10"`
}

type UpdateRoomInput struct {
	Name     string `json:"name" binding:"required" example:"Meeting Room 1"`
	Capacity int    `json:"capacity" binding:"required" example:"10"`
}

type RoomResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
