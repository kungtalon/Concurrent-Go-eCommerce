package main

import (
	"context"
	"google.golang.org/grpc"
	pb "jzmall/lightning/proto"
	"log"
	"net"
	"sync"
)

var sum int64 = 0

// count of product to be sold
var productNum int64 = 10000

// mutual exclusive lock
var mutex sync.Mutex

type server struct {
	pb.UnimplementedCheckRemainsServer
}

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	// check whether the count of product has exceeded storage
	if sum < productNum {
		sum += 1
		return true
	}
	return false
}

// TryGetOne is an implementation of the predefined gRPC interface
func (s *server) TryGetOne(ctx context.Context, in *pb.GetOneRequest) (*pb.GetOneReply, error) {
	log.Printf("ProtoBuf Received: %v", in.GetProductId())
	if GetOneProduct() {
		return &pb.GetOneReply{Remaining: "true"}, nil
	}
	return &pb.GetOneReply{Remaining: "false"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCheckRemainsServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
