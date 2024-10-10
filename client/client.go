package rstore_client

import (
	"context"
	"fmt"

	pb "github.com/brotherlogic/mstore/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MStoreClient interface {
	Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error)
	Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error)
	GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
}

type rClient struct {
	gClient pb.MStoreServiceClient
}

func GetClient() (MStoreClient, error) {
	conn, err := grpc.Dial("mstore.mstore:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024)))
	if err != nil {
		return nil, fmt.Errorf("dial error on %v -> %w", "rmtore.mstore:8080", err)
	}
	return &rClient{gClient: pb.NewMStoreServiceClient(conn)}, nil
}

func (c *rClient) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	return c.gClient.Read(ctx, req)
}

func (c *rClient) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	return c.gClient.Write(ctx, req)
}

func (c *rClient) GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	return c.gClient.GetKeys(ctx, req)
}

func (c *rClient) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return c.gClient.Delete(ctx, req)
}
