package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/servicios/events"
	"github.com/servicios/kafkax"
	"github.com/servicios/repository"
)

func idStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

func consumeRoomEvents(ctx context.Context, brokers []string, repo repository.ServicioRepository) {
	handle := func(ctx context.Context, msg kafka.Message) error {
		// Filtro por tipo
		if kafkax.EventType(msg) != events.TypeRoomDeleted {
			return nil
		}

		var evt events.RoomDeleted
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			return nil
		}

		// Buscar servicios de la habitación
		servicios, err := repo.GetByRoomID(idStr(evt.RoomID))
		if err != nil {
			return err
		}

		for _, s := range servicios {
			if err := repo.Delete(idStr(s.ID)); err != nil {
				return err
			}
			log.Printf("[kafka] habitación %d borrada -> servicio %d borrado (partition=%d offset=%d)",
				evt.RoomID, s.ID, msg.Partition, msg.Offset)
		}
		return nil
	}
	log.Printf("[kafka] servicios escuchando %s (grupo %s)", events.TopicRooms, events.GroupServicios)
	if err := kafkax.Consume(ctx, brokers, events.TopicRooms, events.GroupServicios, handle); err != nil {
		log.Printf("[kafka] consumidor detenido: %v", err)
	}
}
