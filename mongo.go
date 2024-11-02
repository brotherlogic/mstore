package main

import (
	"context"
	"log"

	pb "github.com/brotherlogic/mstore/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mongoClient struct {
	client *mongo.Client
}

func (m *mongoClient) Init(ctx context.Context) error {
	r := m.client.Database("proto").RunCommand(context.Background(), bson.D{{"createUser", "admin"},
		{"pwd", "pass"}, {"roles", []bson.M{{"role": "root", "db": "admin"}}}})
	if r.Err() != nil {
		return r.Err()
	}

	return r.Err()
}

func (m *mongoClient) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplmented")
}

func (m *mongoClient) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	collection := m.client.Database("proto").Collection("protos")
	_, err := collection.InsertOne(ctx, bson.D{
		{"name", req.GetKey()},
		{"value", string(req.GetValue().GetValue())}})
	log.Printf("Write: %v", err)
	return &pb.WriteResponse{}, err
}

func (m *mongoClient) GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplmented")
}

func (m *mongoClient) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplmented")
}

func (m *mongoClient) Count(ctx context.Context, req *pb.CountRequest) (*pb.CountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplemented")
}
