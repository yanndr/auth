package server

import (
	pb "auth/api/v1"
	"auth/internal/model"
	"auth/internal/services"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.AuthServer
	userService services.UserService
}

func NewServer(userService services.UserService) *AuthServer {
	return &AuthServer{
		userService: userService,
	}
}

func (a *AuthServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	err := a.userService.Create(ctx, model.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	return &pb.CreateUserResponse{Success: true}, nil
}

func (a *AuthServer) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	token, err := a.userService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", err)
	}

	return &pb.AuthenticateResponse{
		Success: true,
		Token:   token,
	}, nil
}
