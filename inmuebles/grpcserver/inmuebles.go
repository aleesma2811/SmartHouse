package grpcserver

import (
	"context"
	"errors"
	"strconv"

	"github.com/inmuebles/models"
	"github.com/inmuebles/pb"
	"github.com/inmuebles/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type InmuebleServer struct {
	pb.UnimplementedInmuebleServiceServer
	repo repository.InmuebleRepository
}

func NewInmuebleServer(repo repository.InmuebleRepository) *InmuebleServer {
	return &InmuebleServer{repo: repo}
}

func idStr(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}

func toPb(i *models.Inmueble) *pb.Inmueble {
	return &pb.Inmueble{
		Id:        uint32(i.ID),
		Nombre:    i.Nombre,
		Direccion: i.Direccion,
		Ciudad:    i.Ciudad,
		Tipo:      i.Tipo,
	}
}

func (s *InmuebleServer) GetInmueble(ctx context.Context, req *pb.GetInmuebleRequest) (*pb.Inmueble, error) {
	inmueble, err := s.repo.GetByID(idStr(req.GetId()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "inmueble %d no existe", req.GetId())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toPb(inmueble), nil
}

func (s *InmuebleServer) ListInmuebles(ctx context.Context, req *pb.ListInmueblesRequest) (*pb.ListInmueblesResponse, error) {
	inmuebles, err := s.repo.GetAll()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := make([]*pb.Inmueble, 0, len(inmuebles))
	for i := range inmuebles {
		out = append(out, toPb(&inmuebles[i]))
	}
	return &pb.ListInmueblesResponse{Inmuebles: out}, nil
}

func (s *InmuebleServer) InmuebleExists(ctx context.Context, req *pb.GetInmuebleRequest) (*pb.ExistsResponse, error) {
	_, err := s.repo.GetByID(idStr(req.GetId()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ExistsResponse{Exists: false}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ExistsResponse{Exists: true}, nil
}
