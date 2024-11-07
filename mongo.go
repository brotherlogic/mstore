package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/brotherlogic/mstore/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type mongoClient struct {
	client *mongo.Client
}

func (m *mongoClient) Init(ctx context.Context) error {

	r := m.client.Database("proto").RunCommand(context.Background(), bson.D{{"createUser", "user"},
		{"pwd", "pass"}, {"roles", []bson.M{{"role": "readWrite", "db": "proto"}}}})
	if r.Err() != nil {
		return r.Err()
	}

	return r.Err()
}

func (m *mongoClient) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	var result struct {
		Value string
	}

	collection := m.client.Database("proto").Collection("protos")
	err := collection.FindOne(ctx, bson.D{{"name", req.GetKey()}}).Decode(&result)
	log.Printf("Read: %v", err)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, status.Errorf(codes.NotFound, "unable to locate %v", req.GetKey())
	} else if err != nil {
		return nil, fmt.Errorf("error in mongo read: %v", err)
	}

	return &pb.ReadResponse{Value: &anypb.Any{Value: []byte(result.Value)}}, nil
}

func (m *mongoClient) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	collection := m.client.Database("proto").Collection("protos")

	// Pre clear this key temporarily whislt we deal with the issue of writes
	_, err := collection.DeleteMany(ctx, bson.D{{"name", req.GetKey()}})
	if err != nil {
		log.Printf("Unable to delete on write path: %v", err)
		return nil, err
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(
		ctx,
		bson.D{{"name", req.GetKey()}},
		bson.D{
			{"name", req.GetKey()},
			{"value", string(req.GetValue().GetValue())}},
		opts)
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
