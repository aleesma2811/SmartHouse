package repository

import (
	"github.com/back/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	GetAll() ([]models.Room, error)
	GetByID(id string) (*models.Room, error)
	Create(room *models.Room) error
	Delete(id string) error
}

// Esta estructura es la que habla con Postgres
type gormRoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &gormRoomRepository{db: db}
}

// Receiver (this/self de gormRoomRepository)
func (repo *gormRoomRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	err := repo.db.Find(&rooms).Error
	return rooms, err
}

func (repo *gormRoomRepository) GetByID(id string) (*models.Room, error) {
	var room models.Room
	err := repo.db.First(&room, id).Error // Trae un solo registro filtrado por el ID

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (repo *gormRoomRepository) Create(room *models.Room) error {
	return repo.db.Create(room).Error
}

func (repo *gormRoomRepository) Delete(id string) error {
	return repo.db.Unscoped().Delete(&models.Room{}, id).Error // Struct vacía, no carga los datos del room, únicamente el ID
}
