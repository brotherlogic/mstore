package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	ghbpb "github.com/brotherlogic/githubridge/proto"
	pb "github.com/brotherlogic/mstore/proto"

	ghbclient "github.com/brotherlogic/githubridge/client"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	port         = flag.Int("port", 8080, "The server port.")
	metricsPort  = flag.Int("metrics_port", 8081, "Metrics port")
	mongoAddress = flag.String("mongo", "mongodb://localhost:27017", "Connection String")
)

var (
	wCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mstore_wcount",
	}, []string{"client", "code"})
)

type Server struct {
	rdb     *redis.Client
	gclient ghbclient.GithubridgeClient

	mongoClient *mongoClient
}

type mstore interface {
	Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error)
	Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error)
	GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
}

func (s *Server) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	return s.mongoClient.Read(ctx, req)
}

func (s *Server) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	return s.mongoClient.Write(ctx, req)
}

func (s *Server) GetKeys(ctx context.Context, req *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	return s.mongoClient.GetKeys(ctx, req)
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return s.mongoClient.Delete(ctx, req)
}

func (s *Server) Count(ctx context.Context, req *pb.CountRequest) (*pb.CountResponse, error) {
	return s.mongoClient.Count(ctx, req)
}

func main() {
	flag.Parse()

	s := &Server{}
	client, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to reach GHB")
	}
	s.gclient = client

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	clientOpts := options.Client().ApplyURI(*mongoAddress)
	mclient, err := mongo.Connect(ctx, clientOpts)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err = mclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	s.mongoClient = &mongoClient{client: mclient}

	// Check that we've created a user for the db
	err = s.mongoClient.Init(ctx)
	log.Printf("Init err: %v", err)

	err = mclient.Ping(ctx, readpref.Primary())
	if err != nil {
		_, err = s.gclient.CreateIssue(ctx, &ghbpb.CreateIssueRequest{
			User:  "brotherlogic",
			Repo:  "rstore",
			Title: "Mongo Ping Failure",
			Body:  fmt.Sprintf("Error: %v", err),
		})
		if err != nil {
			panic(err)
		}
	}
	log.Printf("PING: %v", err)
	cancel()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("mstore failed to listen on the serving port %v: %v", *port, err)
	}
	size := 1024 * 1024 * 1000
	gs := grpc.NewServer(
		grpc.MaxSendMsgSize(size),
		grpc.MaxRecvMsgSize(size),
	)
	pb.RegisterMStoreServiceServer(gs, s)
	log.Printf("mstore is listening on %v", lis.Addr())

	// Setup prometheus export
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
	}()

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("mstore failed to serve: %v", err)
	}
}
