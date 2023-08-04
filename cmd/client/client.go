package main

import (
	"auth/pkg/model"
	"auth/pkg/pb"
	"context"
	ctls "crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
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

		tlsConfig, err := SetupTLSConfig(model.TLSConfig{
			CertFile:      "cert/client_cert.pem",
			KeyFile:       "cert/client_key.pem",
			CAFile:        "cert/ca_cert.pem",
			ServerAddress: "localhost:50051",
		})
		if err != nil {
			log.Fatal(err)
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	// Create a connection with the TLS credentials
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("could not dial %s: %s", *serverAddr, err)
	}

	// Initialize the client and make the request
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

func SetupTLSConfig(cfg model.TLSConfig) (*ctls.Config, error) {
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
