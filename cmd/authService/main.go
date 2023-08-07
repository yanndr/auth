package main

import (
	"auth/pkg/config"
	"auth/pkg/jwt"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/stores"
	"auth/pkg/stores/pg"
	"auth/pkg/validators"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net"

	"log"
)

var Version = "v0.1-dev"

func main() {
	// setup minimal log for this app
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(logger)

	logger.Info("starting auth service", zap.String("Version", Version))

	configuration, err := config.LoadConfiguration("config", "./config")
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.Open(configuration.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set all the dependencies
	userValidator := validators.NewUserValidator(validator.New(), validators.NewPasswordValidator(configuration.Password))
	jwtGenerator := jwt.NewTokenGenerator(configuration.Token)
	userStore := stores.NewPgUserStore(pg.New(db))
	userService := services.NewUserService(userStore, userValidator, 10)
	authService := services.NewJwtAuthService(userStore, jwtGenerator)
	grpcAuthServer := server.NewAuthServer(userService, authService)

	srv, err := server.NewGrpcServer(configuration.TLSConfig, grpcAuthServer)

	if err != nil {
		log.Fatal(err)
	}

	logger.Info(
		"service started",
		zap.String("Network", configuration.Network),
		zap.String("Address", configuration.Address),
		zap.Int("Port", configuration.GRPCPort),
	)

	lis, err := net.Listen(configuration.Network, fmt.Sprintf("%s:%v", configuration.Address, configuration.GRPCPort))
	if err != nil {
		log.Fatalf("could not listen on %s:%v: %s", configuration.Address, configuration.GRPCPort, err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("grpc serve error: %s", err)
	}
}
