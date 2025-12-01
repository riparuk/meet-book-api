package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/model"
	"github.com/riparuk/meet-book-api/internal/repository"
)

type BookingHandler struct {
	repo repository.BookingRepository
}

func NewBookingHandler(repo repository.BookingRepository) *BookingHandler {
	return &BookingHandler{repo}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new room booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param input body model.CreateBookingInput true "Booking details"
// @Success 201 {object} model.BookingResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var input model.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if room is available
	available, err := h.repo.IsRoomAvailable(input.RoomID, input.StartTime, input.EndTime, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check room availability"})
		return
	}
	if !available {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room is not available for the selected time slot"})
		return
	}

	booking := model.Booking{
		RoomID:    input.RoomID,
		UserID:    input.UserID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    model.BookingStatusActive,
	}

	if err := h.repo.Create(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking.ToResponse()})
}

// GetBooking godoc
// @Summary Get a booking by ID
// @Description Get a booking by its ID
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} model.BookingResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	booking, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch booking"})
		return
	}

	if booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking.ToResponse()})
}

// GetUserBookings godoc
// @Summary Get all bookings for a user
// @Description Get all bookings for a specific user
// @Tags bookings
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} model.BookingResponse
// @Router /bookings/users/{user_id} [get]
func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	bookings, err := h.repo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user bookings"})
		return
	}

	responses := make([]model.BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// UpdateBooking godoc
// @Summary Update a booking
// @Description Update an existing booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param input body model.UpdateBookingInput true "Booking update details"
// @Success 200 {object} model.BookingResponse
// @Router /bookings/{id} [put]
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	// Get existing booking
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch booking"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	var input model.UpdateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Status != nil {
		existing.Status = *input.Status
	}
	if input.StartTime != nil {
		existing.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		existing.EndTime = *input.EndTime
	}

	// If time is being updated, check room availability
	if input.StartTime != nil || input.EndTime != nil {
		available, err := h.repo.IsRoomAvailable(existing.RoomID, existing.StartTime, existing.EndTime, &existing.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check room availability"})
			return
		}
		if !available {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room is not available for the selected time slot"})
			return
		}
	}

	if err := h.repo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": existing.ToResponse()})
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking
// @Tags bookings
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} model.BookingResponse
// @Router /bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	// Get existing booking
	existing, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch booking"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	if existing.Status == model.BookingStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking is already cancelled"})
		return
	}

	if err := h.repo.Cancel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel booking"})
		return
	}

	existing.Status = model.BookingStatusCancelled
	c.JSON(http.StatusOK, gin.H{"data": existing.ToResponse()})
}

// GetUpcomingBookings godoc
// @Summary Get upcoming bookings
// @Description Get a list of all upcoming bookings
// @Tags bookings
// @Produce json
// @Success 200 {array} model.BookingResponse
// @Router /bookings/upcoming [get]
func (h *BookingHandler) GetUpcomingBookings(c *gin.Context) {
	bookings, err := h.repo.GetUpcomingBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming bookings"})
		return
	}

	responses := make([]model.BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetRoomBookings godoc
// @Summary Get bookings for a specific room
// @Description Get a list of all bookings for a specific room
// @Tags bookings
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {array} model.BookingResponse
// @Router /bookings/room/{room_id} [get]
func (h *BookingHandler) GetRoomBookings(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	bookings, err := h.repo.FindByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch room bookings"})
		return
	}

	responses := make([]model.BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetRoomBookingsByDate godoc
// @Summary Get bookings for a specific room on a specific date
// @Description Get a list of all bookings for a specific room on a specific date
// @Tags bookings
// @Produce json
// @Param room_id path string true "Room ID"
// @Param date path string true "Date"
// @Success 200 {array} model.BookingResponse
// @Router /bookings/room/{room_id}/{date} [get]
func (h *BookingHandler) GetRoomBookingsByDate(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	date, err := time.Parse("2006-01-02", c.Param("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	bookings, err := h.repo.FindByRoomIDAndDate(roomID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch room bookings"})
		return
	}

	responses := make([]model.BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = b.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
