package main

import (
	"context"
	"log"
	"os"
	"time"

	pbrs "github.com/brotherlogic/mstore/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := grpc.NewClient(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Bad dial: %v", err)
	}

	client := pbrs.NewMStoreServiceClient(conn)
	pval := &pbrs.CountRequest{Counter: "TestingCounter"}
	res, err := proto.Marshal(pval)
	if err != nil {
		log.Fatalf("Bad marshal: %v", err)
	}
	_, err = client.Write(ctx, &pbrs.WriteRequest{Key: "test/counter", Value: &anypb.Any{Value: res}})
	if err != nil {
		log.Fatalf("Bad write: %v", err)
	}

	pval = &pbrs.CountRequest{Counter: "TestinmgCounter2"}
	res, err = proto.Marshal(pval)
	if err != nil {
		log.Fatalf("Bad marshal: %v", err)
	}
	_, err = client.Write(ctx, &pbrs.WriteRequest{Key: "test/counter", Value: &anypb.Any{Value: res}})
	if err != nil {
		log.Fatalf("Bad write: %v", err)
	}

	ret, err := client.Read(ctx, &pbrs.ReadRequest{Key: "test/counter"})
	if err != nil {
		log.Fatalf("Bad read: %v", err)
	}

	npval := &pbrs.CountRequest{}
	err = proto.Unmarshal(ret.GetValue().GetValue(), npval)
	if err != nil {
		log.Fatalf("Bad unmarshal: %v", err)
	}

	log.Printf("%v (%v) vs %v (%v)", pval, npval, res, ret.GetValue().GetValue())

}
