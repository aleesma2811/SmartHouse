package repository

import (
	"github.com/inmuebles/models"
	"gorm.io/gorm"
)

type InmuebleRepository interface {
	GetAll() ([]models.Inmueble, error)
	GetByID(id string) (*models.Inmueble, error)
	Create(inmueble *models.Inmueble) error
	Update(inmueble *models.Inmueble) error
	Delete(id string) error
}

type gormInmuebleRepository struct {
	db *gorm.DB
}

func NewInmuebleRepository(db *gorm.DB) InmuebleRepository {
	return &gormInmuebleRepository{db: db}
}

func (repo *gormInmuebleRepository) GetAll() ([]models.Inmueble, error) {
	var inmuebles []models.Inmueble
	err := repo.db.Find(&inmuebles).Error
	return inmuebles, err
}

func (repo *gormInmuebleRepository) GetByID(id string) (*models.Inmueble, error) {
	var inmueble models.Inmueble
	err := repo.db.First(&inmueble, id).Error

	if err != nil {
		return nil, err
	}
	return &inmueble, nil
}

func (repo *gormInmuebleRepository) Create(inmueble *models.Inmueble) error {
	return repo.db.Create(inmueble).Error
}

func (repo *gormInmuebleRepository) Update(inmueble *models.Inmueble) error {
	return repo.db.Save(inmueble).Error
}

func (repo *gormInmuebleRepository) Delete(id string) error {
	return repo.db.Unscoped().Delete(&models.Inmueble{}, id).Error
}
