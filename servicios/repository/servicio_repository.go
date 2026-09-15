package repository

import (
	"github.com/servicios/models"
	"gorm.io/gorm"
)

type ServicioRepository interface {
	GetAll() ([]models.Servicio, error)
	GetByRoomID(roomID string) ([]models.Servicio, error)
	GetByID(id string) (*models.Servicio, error)
	Create(plug *models.Servicio) error
	Update(plug *models.Servicio) error
	Delete(id string) error
}

type gormServicioRepository struct {
	db *gorm.DB
}

func NewServicioRepository(db *gorm.DB) ServicioRepository {
	return &gormServicioRepository{db: db}
}

func (repo *gormServicioRepository) GetAll() ([]models.Servicio, error) {
	var servicios []models.Servicio
	err := repo.db.Find(&servicios).Error
	return servicios, err
}

func (repo *gormServicioRepository) GetByRoomID(roomID string) ([]models.Servicio, error) {
	var servicios []models.Servicio
	err := repo.db.Where("room_id = ?", roomID).Find(&servicios).Error
	return servicios, err
}

func (repo *gormServicioRepository) GetByID(id string) (*models.Servicio, error) {
	var servicio models.Servicio
	err := repo.db.First(&servicio, id).Error

	if err != nil {
		return nil, err
	}
	return &servicio, nil
}

func (repo *gormServicioRepository) Create(servicio *models.Servicio) error {
	return repo.db.Create(servicio).Error
}

func (repo *gormServicioRepository) Update(servicio *models.Servicio) error {
	return repo.db.Save(servicio).Error
}

func (repo *gormServicioRepository) Delete(id string) error {
	return repo.db.Unscoped().Delete(&models.Servicio{}, id).Error
}
