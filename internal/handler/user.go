package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/model"
	"github.com/riparuk/meet-book-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo    repository.UserRepository
	bookingRepo repository.BookingRepository
}

func NewUserHandler(userRepo repository.UserRepository, bookingRepo repository.BookingRepository) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		bookingRepo: bookingRepo,
	}
}

// CreateMyBooking godoc
// @Summary Create a new booking for the authenticated user
// @Description Create a new room booking for the currently authenticated user
// @Tags me
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.CreateMyBookingInput true "Booking details"
// @Success 201 {object} model.BookingResponse
// @Router /me/bookings [post]
func (h *UserHandler) CreateMyBooking(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var input model.CreateMyBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create booking object
	booking := model.Booking{
		RoomID:    input.RoomID,
		UserID:    userUUID,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Status:    model.BookingStatusActive,
	}

	// Validate booking
	if err := booking.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if room is available
	available, err := h.bookingRepo.IsRoomAvailable(input.RoomID, input.StartTime, input.EndTime, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check room availability"})
		return
	}
	if !available {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room is not available for the selected time slot"})
		return
	}

	// Create the booking
	if err := h.bookingRepo.Create(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	// Get the created booking with related data
	createdBooking, err := h.bookingRepo.FindByID(booking.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdBooking.ToResponse()})
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param name body model.CreateUserInput true "name"
// @Success 201 {object} model.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input model.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := h.userRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// Profile godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile
// @Tags me
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Router /me [get]
func (h *UserHandler) Profile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userRepo.FindByID(userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
