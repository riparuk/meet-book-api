package repository

import (
	"github.com/google/uuid"
	"github.com/riparuk/meet-book-api/internal/model"
	"gorm.io/gorm"
)

type RoomRepository interface {
	FindAll() ([]model.Room, error)
	Create(room *model.Room) error
	FindByID(id uuid.UUID) (*model.Room, error)
	Update(room *model.Room) error
	Delete(id uuid.UUID) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) FindAll() ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) FindByID(id uuid.UUID) (*model.Room, error) {
	var room model.Room
	err := r.db.First(&room, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Room{}, "id = ?", id).Error
}
