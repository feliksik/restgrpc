package main

import (
	"log"
	"net"

	pb "github.com/feliksik/restgrpc/example/compute"
	"github.com/golang/go/src/fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Compute(ctx context.Context, in *pb.ComputeRequest) (*pb.ComputeResponse, error) {
	res := in.A*in.B + in.C
	str := fmt.Sprintf("hello there!! %d * %d + %d = %d", in.A, in.B, in.C, res)
	log.Println("serving " + str)
	return &pb.ComputeResponse{
		Result: res,
		Ser:    str,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterComputeServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
