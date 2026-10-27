package grpcserver

import (
	"context"
	"errors"
	"strconv"

	"github.com/back/models"
	"github.com/back/pb"
	"github.com/back/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RoomServer struct {
	pb.UnimplementedRoomServiceServer
	repo repository.RoomRepository
}

// Constructor con inyección de dependencias
func NewRoomServer(repo repository.RoomRepository) *RoomServer {
	return &RoomServer{repo: repo}
}

// Convertir ID de string a uint32
func idStr(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}

// Traducr de GORM a protobuf
func toPb(r *models.Room) *pb.Room {
	return &pb.Room{
		Id:         uint32(r.ID),
		Name:       r.Name,
		InmuebleId: uint32(r.InmuebleID),
	}
}

func (s *RoomServer) GetRoom(ctx context.Context, req *pb.GetRoomRequest) (*pb.Room, error) {
	room, err := s.repo.GetByID(idStr(req.GetId()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "room %d no existe", req.GetId())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toPb(room), nil
}

func (s *RoomServer) ListRooms(ctx context.Context, req *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	rooms, err := s.repo.GetAll() // Llama a GetAll()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := make([]*pb.Room, 0, len(rooms))
	for i := range rooms {
		if req.GetInmuebleId() != 0 && uint32(rooms[i].InmuebleID) != req.GetInmuebleId() {
			continue
		}
		out = append(out, toPb(&rooms[i]))
	}
	return &pb.ListRoomsResponse{Rooms: out}, nil
}

func (s *RoomServer) RoomExists(ctx context.Context, req *pb.GetRoomRequest) (*pb.ExistsResponse, error) {
	_, err := s.repo.GetByID(idStr(req.GetId()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ExistsResponse{Exists: false}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ExistsResponse{Exists: true}, nil
}
