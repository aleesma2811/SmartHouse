package models

import "gorm.io/gorm"

const (
	TipoCasa         = "casa"
	TipoDepartamento = "departamento"
)

type Inmueble struct {
	gorm.Model
	Nombre    string
	Direccion string
	Ciudad    string
	Tipo      string
}
