package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/back/events"
	"github.com/back/kafkax"
	"github.com/back/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func idStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func consumeInmueblesEvents(ctx context.Context, brokers []string, repo repository.RoomRepository, roomEvents *kafkax.Publisher) {
	handle := func(ctx context.Context, msg kafka.Message) error {
		// Filtro por tipo
		if kafkax.EventType(msg) != events.TypeInmuebleDeleted {
			return nil
		}

		var evt events.InmuebleDeleted
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			return nil
		}

		// Buscar habitaciones del inmueble
		rooms, err := repo.GetByInmuebleID(idStr(evt.InmuebleID))
		if err != nil {
			return err
		}

		for _, room := range rooms {
			roomEvt := events.RoomDeleted{
				EventID:    uuid.NewString(),
				RoomID:     room.ID,
				InmuebleID: room.InmuebleID,
				DeletedAt:  time.Now().UTC(),
			}
			if err := roomEvents.Publish(ctx, idStr(room.InmuebleID), events.TypeRoomDeleted, roomEvt); err != nil {
				return err
			}
			if err := repo.Delete(idStr(room.ID)); err != nil {
				return err
			}
			log.Printf("[kafka] inmueble %d borrado -> habitación %d borrada (partition=%d offset=%d)",
				evt.InmuebleID, room.ID, msg.Partition, msg.Offset)
		}
		return nil
	}
	log.Printf("[kafka] back escuchando %s (grupo %s)", events.TopicInmuebles, events.GroupBack)
	if err := kafkax.Consume(ctx, brokers, events.TopicInmuebles, events.GroupBack, handle); err != nil {
		log.Printf("[kafka] consumidor detenido: %v", err)
	}
}
