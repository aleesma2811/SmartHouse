package events

import "time"

const (
	// Topic por servicio
	TopicInmuebles = "inmuebles.events"
	TopicRooms     = "rooms.events"
	// Event-type de cada mensaje (headers del mensaje)
	TypeInmuebleDeleted = "inmueble.deleted"
	TypeRoomDeleted     = "room.deleted"
	// Consumer groups
	GroupBack      = "back-service"
	GroupServicios = "servicios-service"
)

type InmuebleDeleted struct {
	EventID    string    `json:"event_id"`
	InmuebleID uint      `json:"inmueble_id"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type RoomDeleted struct {
	EventID    string    `json:"event_id"`
	RoomID     uint      `json:"room_id"`
	InmuebleID uint      `json:"inmueble_id"`
	DeletedAt  time.Time `json:"deleted_at"`
}
