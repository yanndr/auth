package main

import (
	"auth/pkg/config"
	"auth/pkg/jwt"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/stores"
	"auth/pkg/stores/pg"
	"auth/pkg/stores/sqlite"
	"auth/pkg/validators"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net"
)

var Version = "v0.1-dev"

var (
	configFile = pflag.StringP("config_file", "c", "config", "The name of the config file")
	configPath = pflag.StringP("config_path", "p", "./config", "The path to the config file")
	tls        = pflag.Bool("tls", false, "use tls")
)

func main() {
	pflag.Parse()

	// setup log for this app
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(logger)

	err = viper.BindPFlag("tlsconfig.usetls", pflag.CommandLine.Lookup("tls"))
	if err != nil {
		logger.Fatal("cannot map flag to config", zap.Error(err))
	}

	logger.Info("starting auth service", zap.String("Version", Version))

	configuration, err := config.LoadConfiguration(*configFile, *configPath)
	if err != nil {
		logger.Fatal("error reading configuration", zap.Error(err))
	}

	var db *sql.DB
	switch configuration.Database.Type {
	case "sqlite":
		db, err = sqlite.Open(configuration.Database.Path)
	case "postgres":
		db, err = pg.Open(configuration.Database)
	default:
		logger.Fatal("unknown database type")
	}

	if err != nil {
		logger.Fatal("error opening database", zap.Error(err))
	}
	defer db.Close()

	// Set all the dependencies
	userValidator := validators.NewUserValidator(validator.New(), validators.NewPasswordValidator(configuration.Password))
	jwtGenerator := jwt.NewTokenGenerator(configuration.Token)
	userStore := stores.NewPgUserStore(pg.New(db))
	userService := services.NewUserService(userStore, userValidator, 10)
	authService := services.NewJwtAuthService(userStore, jwtGenerator)

	srv, err := server.NewGrpcServer(configuration.TLSConfig, userService, authService)

	if err != nil {
		logger.Fatal("error creating the grpc server", zap.Error(err))
	}

	logger.Info(
		"service started",
		zap.String("Network", configuration.Network),
		zap.String("Address", configuration.Address),
		zap.Int("Port", configuration.GRPCPort),
		zap.Bool("TLS", configuration.TLSConfig.UseTLS),
	)

	lis, err := net.Listen(configuration.Network, fmt.Sprintf("%s:%v", configuration.Address, configuration.GRPCPort))
	if err != nil {
		logger.Fatal(
			"could not start listener",
			zap.String("address", configuration.Address),
			zap.Int("port", configuration.GRPCPort),
			zap.Error(err))
	}

	if err := srv.Serve(lis); err != nil {
		logger.Fatal("grpc server error", zap.Error(err))
	}
}
