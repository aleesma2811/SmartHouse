package repository

import (
	"github.com/servicios/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServicioRepository interface {
	GetAll() ([]models.Servicio, error)
	GetByRoomID(roomID string) ([]models.Servicio, error)
	GetByID(id string) (*models.Servicio, error)
	Create(plug *models.Servicio) error
	Update(plug *models.Servicio) error
	Delete(id string) error

	// Consumo por streaming
	AddLectura(l *models.Lectura) (bool, error)
	GetLecturas(servicioID string, limit int) ([]models.Lectura, error)
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
	return repo.db.Select("*").Omit("ConsumoTotal").Save(servicio).Error
}

// Borrar el servicio y todo su historial
func (repo *gormServicioRepository) Delete(id string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("servicio_id = ?", id).Delete(&models.Lectura{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&models.Servicio{}, id).Error
	})
}

func (repo *gormServicioRepository) AddLectura(l *models.Lectura) (bool, error) {
	added := false
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		// ON CONFLICT DO NOTHING (si event_id ya existe, no inserta nada)
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(l)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		// UPDATE ATÓMICO (se suma el dato sin leer el valor viejo)
		upd := tx.Model(&models.Servicio{}).
			Where("id = ?", l.ServicioID).
			UpdateColumn("consumo_total", gorm.Expr("consumo_total + ?", l.Valor))

		if upd.Error != nil {
			return upd.Error
		}
		if upd.RowsAffected == 0 {
			return gorm.ErrRecordNotFound // Servicio borrado
		}
		added = true
		return nil
	})
	return added, err
}

func (repo *gormServicioRepository) GetLecturas(servicioID string, limit int) ([]models.Lectura, error) {
	var lecturas []models.Lectura
	err := repo.db.Where("servicio_id = ?", servicioID).
		Order("timestamp desc").
		Limit(limit).
		Find(&lecturas).Error

	if err != nil {
		return nil, err
	}
	for i, j := 0, len(lecturas)-1; i < j; i, j = i+1, j-1 {
		lecturas[i], lecturas[j] = lecturas[j], lecturas[i]
	}
	return lecturas, nil
}
