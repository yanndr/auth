package main

import (
	pb "auth/api/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	serverAddr = "localhost:50051"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: "boby", Password: "gloups"})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Success)

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tResp, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Username: "boby", Password: "gloups"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tResp.Token)
}
