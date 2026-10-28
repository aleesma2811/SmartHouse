package models

import "gorm.io/gorm"

// Tipos de servicios
const (
	TipoLuz  = "luz"
	TipoAgua = "agua"
	TipoGas  = "gas"
)

type Servicio struct {
	gorm.Model
	Name string
	Tipo string
	// Consumo por hora
	Consumo float32 // kWh (luz), litros (agua), m3 (gas)
	RoomID  uint
	On      bool `gorm:"default:false"`
	// Float64 porque suma varios float32
	ConsumoTotal float64 `gorm:"default:0"`
}
