package main

import (
	"auth/pkg/jwt"
	"auth/pkg/model"
	"auth/pkg/pb"
	"auth/pkg/server"
	"auth/pkg/services"
	"auth/pkg/store"
	"auth/pkg/store/pg"
	"auth/pkg/validators"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"log"
	"net"
	"os"
	"strings"
	"time"
)

var Version = "v0.1-dev"

func main() {
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
		Validator:         validator.New(),
		PasswordValidator: validators.NewPasswordValidator(configuration.Password),
	}
	jwtGenerator := jwt.NewGenerator(jwtgo.SigningMethodES256, configuration.Token.SignedKey, "authService", time.Minute*time.Duration(configuration.Token.ExpDuration))

	s := services.NewUserService(us, jwtGenerator, userValidator)

	var opts []grpc.ServerOption
	if configuration.TLSConfig.UseTLS {
		tlsConfig, err := SetupTLSConfig(configuration.TLSConfig)
		if err != nil {
			log.Fatal(err)
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	lis, err := net.Listen(configuration.Network, fmt.Sprintf("%s:%v", configuration.Address, configuration.GRPCPort))
	if err != nil {
		log.Fatalf("could not list on %s:%v: %s", configuration.Address, configuration.GRPCPort, err)
	}

	srv := grpc.NewServer(opts...)

	pb.RegisterAuthServer(srv, server.NewServer(s))

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("grpc serve error: %s", err)
	}
}

func SetupTLSConfig(cfg model.TLSConfig) (*tls.Config, error) {
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
