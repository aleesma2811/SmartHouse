package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/servicios/db"
	"github.com/servicios/grpcclient"
	"github.com/servicios/grpcserver"
	"github.com/servicios/models"
	"github.com/servicios/pb"
	"github.com/servicios/repository"
	"github.com/servicios/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":5001"

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

	db.DB.AutoMigrate(models.Servicio{})

	// Inyección de dependencias
	servicioRepo := repository.NewServicioRepository(db.DB)
	servicioHandler := routes.NewServicioHandler(servicioRepo)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("no se pudo abrir el puerto gRPC %s: %v", grpcPort, err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterServicioServiceServer(grpcSrv, grpcserver.NewServicioServer(servicioRepo))
	reflection.Register(grpcSrv) // permite `grpcurl list` / `describe` sin el .proto
	go func() {
		log.Printf("gRPC (ServicioService) escuchando en %s", grpcPort)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("servidor gRPC: %v", err)
		}
	}()

	roomClient, err := grpcclient.NewRoomClient(envOr("ROOMS_GRPC_ADDR", "nginx:8001"))
	if err != nil {
		log.Fatalf("cliente gRPC de rooms: %v", err)
	}
	defer roomClient.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/heartbeat", heartbeatHandler).Methods("GET")

	// Plugs routes
	r.HandleFunc("/servicios", servicioHandler.GetServiciosHandler).Methods("GET")
	r.HandleFunc("/servicios/{id}", servicioHandler.GetServicioHandler).Methods("GET")
	// Antes de crear el servicio, valida por gRPC que la room exista
	r.HandleFunc("/servicios", validateRoom(roomClient, servicioHandler.PostServiciosHandler)).Methods("POST")
	r.HandleFunc("/servicios/{id}", servicioHandler.UpdateServicioHandler).Methods("PUT")
	r.HandleFunc("/servicios/{id}", servicioHandler.DeleteServiciosHandler).Methods("DELETE")

	serve(r, "servicios", 4001)
}
