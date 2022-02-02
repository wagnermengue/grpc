package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/wagnermengue/grpc/pb"
    "github.com/wagnermengue/grpc/services"
)

func main() {
    lis, err := net.Listen("tcp", "localhost:50051")
    if err != nil {
        log.Fatalf("Could not connect: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
    reflection.Register(grpcServer)

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Could not server: %v", err)
    }
}