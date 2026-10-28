package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/back/db"
	"github.com/back/events"
	"github.com/back/grpcclient"
	"github.com/back/grpcserver"
	"github.com/back/kafkax"
	"github.com/back/models"
	"github.com/back/pb"
	"github.com/back/repository"
	"github.com/back/routes"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":5000"

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := db.DB.DB()
	if err == nil {
		// Si se obtuvo la conexión, hacer un ping a Posgres con límite de 2s
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		err = sqlDB.PingContext(ctx)
	}
	if err != nil {
		http.Error(w, "DB unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Back alive"))
}

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.Room{})

	// Kafka
	// Back publica en rooms.events y onsume inmuebles.events
	kafkaBroker := envOr("KAFKA_BROKER", "kafka:9092")
	brokers := []string{kafkaBroker}
	for _, topic := range []string{events.TopicInmuebles, events.TopicRooms} {
		if err := kafkax.EnsureTopic(kafkaBroker, topic, 3); err != nil {
			log.Fatalf("kafka: no se pudo crear el tópico %s: %v", topic, err)
		}
	}
	roomEvents := kafkax.NewPublisher(brokers, events.TopicRooms)
	defer roomEvents.Close()

	// Inyección de dependencias
	roomRepo := repository.NewRoomRepository(db.DB)
	roomHandler := routes.NewRoomHandler(roomRepo, roomEvents)

	// Consumidor de inmuebles.events: borra habitaciones y publica room.deleted
	kafkaCtx, stopKafka := context.WithCancel(context.Background())
	defer stopKafka()
	go consumeInmueblesEvents(kafkaCtx, brokers, roomRepo, roomEvents)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("no se pudo abrir el puerto gRPC %s: %v", grpcPort, err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterRoomServiceServer(grpcSrv, grpcserver.NewRoomServer(roomRepo))
	reflection.Register(grpcSrv) // permite `grpcurl list` / `describe` sin el .proto
	go func() {
		log.Printf("gRPC (RoomService) escuchando en %s", grpcPort)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("servidor gRPC: %v", err)
		}
	}()

	inmuebleClient, err := grpcclient.NewInmuebleClient(envOr("INMUEBLES_GRPC_ADDR", "nginx:8001"))
	if err != nil {
		log.Fatalf("cliente gRPC de inmuebles: %v", err)
	}
	defer inmuebleClient.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/heartbeat", heartbeatHandler).Methods("GET")

	// Rooms routes
	r.HandleFunc("/rooms", roomHandler.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms/{id}", roomHandler.GetRoomHandler).Methods("GET")
	// Antes de crear la room, valida por gRPC que el inmueble exista
	r.HandleFunc("/rooms", validateInmueble(inmuebleClient, roomHandler.PostRoomsHandler)).Methods("POST")
	r.HandleFunc("/rooms/{id}", roomHandler.DeleteRoomsHandler).Methods("DELETE")

	// http.ListenAndServe(":4000", r)
	serve(r, "back", 4000)
}
