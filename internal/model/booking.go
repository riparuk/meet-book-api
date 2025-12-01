package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RoomID    uuid.UUID      `json:"room_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	StartTime time.Time      `json:"start_time" gorm:"not null"`
	EndTime   time.Time      `json:"end_time" gorm:"not null"`
	Status    BookingStatus  `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Room Room `json:"room" gorm:"foreignKey:RoomID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type CreateBookingInput struct {
	RoomID    uuid.UUID `json:"room_id" binding:"required"`
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type UpdateBookingInput struct {
	Status    *BookingStatus `json:"status,omitempty"`
	StartTime *time.Time     `json:"start_time,omitempty"`
	EndTime   *time.Time     `json:"end_time,omitempty"`
}

type BookingResponse struct {
	ID        uuid.UUID     `json:"id"`
	RoomID    uuid.UUID     `json:"room_id"`
	UserID    uuid.UUID     `json:"user_id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Status    BookingStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`

	Room Room `json:"room"`
	User User `json:"user"`
}

// ToResponse converts a Booking to a BookingResponse
func (b *Booking) ToResponse() BookingResponse {
	return BookingResponse{
		ID:        b.ID,
		RoomID:    b.RoomID,
		UserID:    b.UserID,
		StartTime: b.StartTime,
		EndTime:   b.EndTime,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Room:      b.Room,
		User:      b.User,
	}
}

// BeforeCreate is a hook that runs before creating a booking
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	b.Status = BookingStatusActive
	return nil
}

// Validate checks if the booking time is valid
func (b *Booking) Validate() error {
	if b.StartTime.IsZero() || b.EndTime.IsZero() {
		return fmt.Errorf("start time and end time are required")
	}

	if b.EndTime.Before(b.StartTime) || b.EndTime.Equal(b.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}
