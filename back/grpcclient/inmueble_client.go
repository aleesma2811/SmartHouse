package grpcclient

import (
	"context"
	"time"

	"github.com/back/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InmuebleClient struct {
	conn   *grpc.ClientConn
	client pb.InmuebleServiceClient
}

func NewInmuebleClient(addr string) (*InmuebleClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &InmuebleClient{conn: conn, client: pb.NewInmuebleServiceClient(conn)}, nil
}

func (c *InmuebleClient) Exists(ctx context.Context, id uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.InmuebleExists(ctx, &pb.GetInmuebleRequest{Id: uint32(id)})
	if err != nil {
		return false, err
	}
	return resp.GetExists(), nil
}

func (c *InmuebleClient) Close() error {
	return c.conn.Close()
}
