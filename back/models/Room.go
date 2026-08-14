package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name  string
	Plugs []Plug `gorm:"foreignKey:RoomID"`
}
