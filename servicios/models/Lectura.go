package models

import "time"

type Lectura struct {
	ID         uint `gorm:"primarykey"`
	ServicioID uint `gorm:"index"`
	Valor      float64
	// Si kafka entrega el mismo evento 2 veces, el segundo INSERT se ignora
	EventID   string    `gorm:"uniqueIndex"`
	TimeStamp time.Time `gorm:"index"`
}
