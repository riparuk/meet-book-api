package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	FindByID(id uuid.UUID) (*model.Booking, error)
	FindByUserID(userID uuid.UUID) ([]model.Booking, error)
	FindByRoomID(roomID uuid.UUID) ([]model.Booking, error)
	FindByRoomIDAndDate(roomID uuid.UUID, date time.Time) ([]model.Booking, error)
	Update(booking *model.Booking) error
	Cancel(id uuid.UUID) error
	IsRoomAvailable(roomID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)
	GetUpcomingBookings() ([]model.Booking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) FindByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.
		Preload("Room").
		Preload("User").
		First(&booking, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) FindByUserID(userID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.
		Preload("Room").
		Where("user_id = ?", userID).
		Order("start_time DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) FindByRoomID(roomID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.
		Where("room_id = ?", roomID).
		Order("start_time DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Cancel(id uuid.UUID) error {
	return r.db.Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", model.BookingStatusCancelled).
		Error
}

func (r *bookingRepository) IsRoomAvailable(roomID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&model.Booking{}).
		Where("room_id = ?", roomID).
		Where("status = ?", model.BookingStatusActive).
		Where("(start_time, end_time) OVERLAPS (?, ?)", startTime, endTime)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count == 0, err
}

func (r *bookingRepository) GetUpcomingBookings() ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.
		Preload("Room").
		Preload("User").
		Where("start_time > ?", time.Now()).
		Where("status = ?", model.BookingStatusActive).
		Order("start_time ASC").
		Find(&bookings).Error
	return bookings, err
}

// FindByRoomIDAndDate returns all bookings for a specific room on a specific date
func (r *bookingRepository) FindByRoomIDAndDate(roomID uuid.UUID, date time.Time) ([]model.Booking, error) {
	var bookings []model.Booking

	// Start of the day (00:00:00)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	// End of the day (23:59:59.999999999)
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	err := r.db.
		Where("room_id = ?", roomID).
		Where("start_time >= ? AND end_time <= ?", startOfDay, endOfDay).
		Order("start_time").
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return bookings, nil
}
