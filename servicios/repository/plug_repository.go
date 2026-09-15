package repository

import (
	"github.com/servicios/models"
	"gorm.io/gorm"
)

type PlugRepository interface {
	GetAll() ([]models.Plug, error)
	GetByRoomID(roomID string) ([]models.Plug, error)
	GetByID(id string) (*models.Plug, error)
	Create(plug *models.Plug) error
	Update(plug *models.Plug) error
	Delete(id string) error
}

type gormPlugRepository struct {
	db *gorm.DB
}

func NewPlugRepository(db *gorm.DB) PlugRepository {
	return &gormPlugRepository{db: db}
}

func (repo *gormPlugRepository) GetAll() ([]models.Plug, error) {
	var plugs []models.Plug
	err := repo.db.Find(&plugs).Error
	return plugs, err
}

func (repo *gormPlugRepository) GetByRoomID(roomID string) ([]models.Plug, error) {
	var plugs []models.Plug
	err := repo.db.Where("room_id = ?", roomID).Find(&plugs).Error
	return plugs, err
}

func (repo *gormPlugRepository) GetByID(id string) (*models.Plug, error) {
	var plug models.Plug
	err := repo.db.First(&plug, id).Error

	if err != nil {
		return nil, err
	}
	return &plug, nil
}

func (repo *gormPlugRepository) Create(plug *models.Plug) error {
	return repo.db.Create(plug).Error
}

func (repo *gormPlugRepository) Update(plug *models.Plug) error {
	return repo.db.Save(plug).Error
}

func (repo *gormPlugRepository) Delete(id string) error {
	return repo.db.Unscoped().Delete(&models.Plug{}, id).Error
}
