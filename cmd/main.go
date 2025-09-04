package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ckalagara/group-a-inventory/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ckalagara/group-a-inventory/proto"
)

func main() {
	log.Printf("Starting gRPC server on port 50052: %v", time.Now())
	ctx := context.Background()
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen with error: %v", err)
	}
	log.Printf("Creating gRPC server: %v", time.Now())
	server := grpc.NewServer()
	reflection.Register(server)
	pb.RegisterServiceServer(server, core.NewService(ctx, "mongodb://mongodb:27017"))
	log.Printf("Serving gRPC server: %v", time.Now())

	err = server.Serve(listener)

	if err != nil {
		log.Fatalf("Failed to serve with error: %v", err)
	}
}
