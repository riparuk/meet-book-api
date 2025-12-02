package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/model"
	"github.com/riparuk/meet-book-api/internal/repository"
)

type RoomHandler struct {
	repo repository.RoomRepository
}

func NewRoomHandler(repo repository.RoomRepository) *RoomHandler {
	return &RoomHandler{repo}
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new meeting room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.CreateRoomInput true "Room details"
// @Success 201 {object} model.Room
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var input model.CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := model.Room{
		Name:     input.Name,
		Capacity: input.Capacity,
	}

	if err := h.repo.Create(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": room})
}

// GetRooms godoc
// @Summary Get all rooms
// @Description Get a list of all meeting rooms
// @Tags rooms
// @Produce json
// @Success 200 {array} model.Room
// @Router /rooms [get]
func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rooms})
}

// GetRoom godoc
// @Summary Get a room by ID
// @Description Get a room by its ID
// @Tags rooms
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} model.Room
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

// UpdateRoom godoc
// @Summary Update a room
// @Description Update an existing room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Param input body model.UpdateRoomInput true "Room details"
// @Success 200 {object} model.Room
// @Router /rooms/{id} [put]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input model.UpdateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	room.Name = input.Name
	room.Capacity = input.Capacity

	if err := h.repo.Update(room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete a room by ID
// @Tags rooms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Room ID"
// @Success 204 "No Content"
// @Router /rooms/{id} [delete]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
