package models

import "gorm.io/gorm"

type Plug struct {
	gorm.Model

	Name       string
	KwhConsump float32
	RoomID     uint
	On         bool `gorm:"default:false"`
}
