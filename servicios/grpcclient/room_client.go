package grpcclient

import (
	"context"
	"time"

	"github.com/servicios/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoomClient struct {
	conn   *grpc.ClientConn
	client pb.RoomServiceClient
}

func NewRoomClient(addr string) (*RoomClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &RoomClient{conn: conn, client: pb.NewRoomServiceClient(conn)}, nil
}

func (c *RoomClient) Exists(ctx context.Context, id uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.RoomExists(ctx, &pb.GetRoomRequest{Id: uint32(id)})
	if err != nil {
		return false, err
	}
	return resp.GetExists(), nil
}

func (c *RoomClient) Close() error {
	return c.conn.Close()
}
