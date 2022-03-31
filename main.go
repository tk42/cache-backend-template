package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/tk42/cache-backend-template/ent"
	"github.com/tk42/cache-backend-template/ent/proto/entpb"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Print("server is starting...")
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		"db", "5432", "postgres", "postgres", "p@ssw0rd"))
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to create schema: %s", err)
	}

	svc := entpb.NewUserService(client)
	server := grpc.NewServer()

	reflection.Register(server) // Failed to list services: server does not support the reflection API
	entpb.RegisterUserServiceServer(server, svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
