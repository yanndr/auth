package main

import (
	"auth/pkg/config"
	"auth/pkg/jwt"
	"auth/pkg/pb"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/stores"
	"auth/pkg/stores/pg"
	"auth/pkg/validators"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"log"
	"net"
	"os"
	"strings"
)

var Version = "v0.1-dev"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(logger)

	logger.Info("starting auth service", zap.String("Version", Version))

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configuration, err := config.LoadConfiguration("config", "./config")
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.Open(configuration.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userStore := stores.NewPgUserStore(pg.New(db))
	userValidator := validators.UserValidator{
		StructValidator:   validator.New(),
		PasswordValidator: validators.NewPasswordValidator(configuration.Password),
	}

	jwtGenerator := jwt.NewGenerator(configuration.Token)

	userService := services.NewUserService(userStore, userValidator, 10)
	authService := &services.JwtAuthService{UserStore: userStore, JwtGenerator: jwtGenerator}

	var opts []grpc.ServerOption
	if configuration.TLSConfig.UseTLS {
		logger.Info("setup TLS")
		tlsConfig, err := setupTLSConfig(configuration.TLSConfig)
		if err != nil {
			log.Fatal(err)
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	lis, err := net.Listen(configuration.Network, fmt.Sprintf("%s:%v", configuration.Address, configuration.GRPCPort))
	if err != nil {
		log.Fatalf("could not listen on %s:%v: %s", configuration.Address, configuration.GRPCPort, err)
	}

	srv := grpc.NewServer(opts...)

	pb.RegisterAuthServer(srv, server.NewServer(userService, authService))
	logger.Info(
		"service started",
		zap.String("Network", configuration.Network),
		zap.String("Address", configuration.Address),
		zap.Int("Port", configuration.GRPCPort),
	)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("grpc serve error: %s", err)
	}
}

func setupTLSConfig(cfg config.TLS) (*tls.Config, error) {
	var err error
	tlsConfig := &tls.Config{}
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(
			cfg.CertFile,
			cfg.KeyFile,
		)
		if err != nil {
			return nil, err
		}
	}
	if cfg.CAFile != "" {
		b, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, err
		}
		ca := x509.NewCertPool()
		ok := ca.AppendCertsFromPEM([]byte(b))
		if !ok {
			return nil, fmt.Errorf(
				"failed to parse root certificate: %q",
				cfg.CAFile,
			)
		}
		tlsConfig.ClientCAs = ca
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ServerName = cfg.ServerAddress
	}
	return tlsConfig, nil
}
