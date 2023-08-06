package server

import (
	"auth/pkg/models"
	"auth/pkg/pb"
	"auth/pkg/services"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"strings"
)

type AuthServer struct {
	pb.AuthServer
	userService services.UserService
	authService services.AuthService
	logger      *zap.Logger
}

func NewServer(userService services.UserService, authService services.AuthService) *AuthServer {
	return &AuthServer{
		userService: userService,
		authService: authService,
		logger:      zap.L().Named("gRPCAuthServer"),
	}
}

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
