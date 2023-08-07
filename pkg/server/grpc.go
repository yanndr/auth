package server

import (
	"auth/pkg/config"
	"auth/pkg/models"
	"auth/pkg/pb"
	"auth/pkg/services"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

type AuthServer struct {
	pb.AuthServer
	userService services.UserService
	authService services.AuthService
	logger      *zap.Logger
}

func NewGrpcServer(configuration config.TLS, authServer pb.AuthServer) (*grpc.Server, error) {
	//Usually get more config here for logging and tracing middleware
	var opts []grpc.ServerOption
	if configuration.UseTLS {
		tlsConfig, err := setupTLSConfig(configuration)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterAuthServer(srv, authServer)

	return srv, nil
}

// NewAuthServer create a new instance of  AuthServer with a services.UserService and a services.AuthService
func NewAuthServer(userService services.UserService, authService services.AuthService) *AuthServer {
	return &AuthServer{
		userService: userService,
		authService: authService,
		logger:      zap.L().Named("gRPCAuthServer"),
	}
}

// CreateUser creates a user form the pb.CreateUserRequest
func (a *AuthServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	a.logger.Info("CreateUser called")
	err := a.userService.Create(ctx, models.User{
		Username: strings.TrimSpace(req.Username),
		Password: strings.TrimSpace(req.Password),
	})

	s, ok := status.FromError(err)
	if err != nil && ok {
		return nil, s.Err()
	} else if err != nil {
		a.logger.Error("unknown error", zap.Error(err))
		return nil, s.Err()
	}

	return &pb.CreateUserResponse{Success: true}, nil
}

// Authenticate a user from the request pb.AuthenticateRequest
func (a *AuthServer) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	a.logger.Info("Authenticate called")
	token, err := a.authService.Authenticate(
		ctx,
		strings.TrimSpace(req.Username),
		strings.TrimSpace(req.Password),
	)
	s, ok := status.FromError(err)
	if err != nil && ok {
		return nil, s.Err()
	} else if err != nil {
		a.logger.Error("unknown error", zap.Error(err))
		return nil, s.Err()
	}

	return &pb.AuthenticateResponse{
		Token: token,
	}, nil
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
