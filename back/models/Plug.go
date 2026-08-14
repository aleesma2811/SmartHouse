package models

import "gorm.io/gorm"

type Plug struct {
	gorm.Model

	Name       string
	KwhConsump float32
	RoomID     uint
	// Room       Room `gorm:"foreignKey:RoomID"`
	On bool `gorm:"default:false"`
}
