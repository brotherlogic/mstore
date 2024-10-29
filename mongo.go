package main

import (
	"context"

	pb "github.com/brotherlogic/mstore/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mongoClient struct {
	client *mongo.Client
}

func (m *mongoClient) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplmented")
}

func (m *mongoClient) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	collection := m.client.Database("proto").Collection("protos")
	_, err := collection.InsertOne(ctx, bson.D{
		{"name", req.GetKey()},
		{"value", string(req.GetValue().GetValue())}})
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
