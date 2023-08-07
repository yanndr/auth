package main

import (
	"auth/pkg/config"
	"auth/pkg/pb"
	"context"
	ctls "crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

var (
//subCmd = flag.NewFlagSet("sub", flag.ExitOnError)

// tls      = pflag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
// username = pflag.String("username", "", "the username")
// password = pflag.String("password", "", "the password")
//
// caFile     = pflag.String("ca_file", "cert/ca_cert.pem", "The file containing the CA root cert file")
// certFile   = pflag.String("cert_file", "cert/client_cert.pem", "The file containing the client cert file")
// keyFile    = pflag.String("key_file", "cert/client_key.pem", "The file containing the client key file")
// serverAddr = pflag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	tls := defaults.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	username := defaults.StringP("username", "u", "", "the username")
	password := defaults.StringP("password", "p", "", "the password")

	caFile := defaults.String("ca_file", "cert/ca_cert.pem", "The file containing the CA root cert file")
	certFile := defaults.String("cert_file", "cert/client_cert.pem", "The file containing the client cert file")
	keyFile := defaults.String("key_file", "cert/client_key.pem", "The file containing the client key file")
	serverAddr := defaults.StringP("addr", "a", "localhost:50051", "The server address in the format of host:port")

	defaults.Parse(os.Args)

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or 'auth' subcommands")
		pflag.PrintDefaults()
		os.Exit(1)
	}

	var cmd func(ctx context.Context, client pb.AuthClient) error
	switch os.Args[1] {
	case "create":
		cmd = func(ctx context.Context, client pb.AuthClient) error {
			response, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: *username, Password: *password})
			if err != nil {
				return err
			}
			fmt.Println(response)
			return nil
		}
	case "auth":
		cmd = func(ctx context.Context, client pb.AuthClient) error {
			response, err := client.Authenticate(ctx, &pb.AuthenticateRequest{Username: *username, Password: *password})
			if err != nil {
				return err
			}
			fmt.Println(response)
			return nil
		}
	default:
		fmt.Println("expected 'create' or 'auth' subcommands")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var opts []grpc.DialOption
	if *tls {
		tlsConfig, err := SetupTLSConfig(config.TLS{
			CertFile:      *certFile,
			KeyFile:       *keyFile,
			CAFile:        *caFile,
			ServerAddress: *serverAddr,
		})
		if err != nil {
			log.Fatal(err)
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("could not dial %s: %s", *serverAddr, err)
	}

	client := pb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = cmd(ctx, client)

	if err != nil {
		log.Fatal(err)
	}
}

func SetupTLSConfig(cfg config.TLS) (*ctls.Config, error) {
	var err error
	tlsConfig := &ctls.Config{}
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		tlsConfig.Certificates = make([]ctls.Certificate, 1)
		tlsConfig.Certificates[0], err = ctls.LoadX509KeyPair(
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

		tlsConfig.RootCAs = ca

		tlsConfig.ServerName = cfg.ServerAddress
	}
	return tlsConfig, nil
}
