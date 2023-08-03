package main

import (
	pb "auth/api/v1"
	"auth/internal/model"
	"auth/internal/services"
	"auth/internal/store"
	"auth/internal/store/pg"
	"auth/server"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "auth_user"
	password = "autPassw@ord"
	dbname   = "auth"

	gRPCPort = ":50051"
)

func main() {

	db, err := pg.Open(host, port, user, password, dbname)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	us := store.NewUserStore(pg.New(db))
	s := services.NewUserService(us)
	err = s.Create(model.User{Username: "bob", Password: "gloups"})
	if err != nil {
		fmt.Println(err)
	}

	t, err := s.Authenticate("yann", "gloups")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(t)

	grpcListener, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	//if *tls {
	//	if *certFile == "" {
	//		*certFile = data.Path("x509/server_cert.pem")
	//	}
	//	if *keyFile == "" {
	//		*keyFile = data.Path("x509/server_key.pem")
	//	}
	//	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	//	if err != nil {
	//		log.Fatalf("Failed to generate credentials: %v", err)
	//	}
	//	opts = []grpc.ServerOption{grpc.Creds(creds)}
	//}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, server.NewServer(s))
	grpcServer.Serve(grpcListener)

}
