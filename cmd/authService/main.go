package main

import (
	"auth/pkg/model"
	"auth/pkg/pb"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/store"
	"auth/pkg/store/pg"
	"auth/pkg/validators"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strings"
)

var Version = "v0.1-dev"

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
)

func main() {
	flag.Parse()
	log.Printf("starting auth service %v", Version)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configuration := &model.Configuration{}
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	db, err := pg.Open(configuration.Database)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	us := store.NewUserStore(pg.New(db))
	userValidator := validators.UserValidator{
		Validator: validator.New(),
		PasswordValidator: validators.PasswordValidator{
			MinLength:  configuration.Password.MinLength,
			MinUpper:   configuration.Password.MinUpperCase,
			MinNumeric: configuration.Password.MinNumeric,
		},
	}
	s := services.NewUserService(us, userValidator)

	grpcListener, err := net.Listen(configuration.Network, fmt.Sprintf("%s:%v", configuration.Address, configuration.Port))
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = "cert/server_cert.pem"
		}
		if *keyFile == "" {
			*keyFile = "cert/server_key.pem"
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, server.NewServer(s))
	log.Printf("listening on %s:%v", configuration.Address, configuration.Port)
	err = grpcServer.Serve(grpcListener)
	if err != nil {
		log.Fatal(err)
	}

}
