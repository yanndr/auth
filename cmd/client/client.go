package main

import (
	"auth/pkg/pb"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "sub.yannd.dev", "The server name used to verify the hostname returned by the TLS handshake")
	username           = flag.String("username", "", "the username")
	password           = flag.String("password", "", "the password")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *tls {
		if *caFile == "" {
			*caFile = "cert/ca_cert.pem"
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: *username, Password: *password})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Success)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tResp, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Username: *username, Password: *password})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tResp.Token)
}
