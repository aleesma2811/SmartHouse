package grpcserver

import (
	"context"
	"errors"
	"strconv"

	"github.com/servicios/models"
	"github.com/servicios/pb"
	"github.com/servicios/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ServicioServer struct {
	pb.UnimplementedServicioServiceServer
	repo repository.ServicioRepository
}

func NewServicioServer(repo repository.ServicioRepository) *ServicioServer {
	return &ServicioServer{repo: repo}
}

func idStr(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}

func toPb(s *models.Servicio) *pb.Servicio {
	return &pb.Servicio{
		Id:      uint32(s.ID),
		Name:    s.Name,
		Tipo:    s.Tipo,
		Consumo: s.Consumo,
		RoomId:  uint32(s.RoomID),
		On:      s.On,
	}
}

func (s *ServicioServer) GetServicio(ctx context.Context, req *pb.GetServicioRequest) (*pb.Servicio, error) {
	servicio, err := s.repo.GetByID(idStr(req.GetId()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "servicio %d no existe", req.GetId())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toPb(servicio), nil
}

func (s *ServicioServer) ListServicios(ctx context.Context, req *pb.ListServiciosRequest) (*pb.ListServiciosResponse, error) {
	var (
		servicios []models.Servicio
		err       error
	)
	if req.GetRoomId() != 0 {
		servicios, err = s.repo.GetByRoomID(idStr(req.GetRoomId()))
	} else {
		servicios, err = s.repo.GetAll()
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := make([]*pb.Servicio, 0, len(servicios))
	for i := range servicios {
		out = append(out, toPb(&servicios[i]))
	}
	return &pb.ListServiciosResponse{Servicios: out}, nil
}
