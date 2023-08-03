package server

import (
	"auth/pkg/model"
	"auth/pkg/pb"
	"auth/pkg/services"
	"auth/pkg/validators"
	"context"
	"errors"
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
		if errors.As(err, &services.UsernameAlreadyExistErr{}) {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
		}
		if errors.As(err, &validators.ValidationErr{}) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		}
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &pb.CreateUserResponse{Success: true}, nil
}

func (a *AuthServer) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	token, err := a.userService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.AutenticationErr) {
			return &pb.AuthenticateResponse{
				Success: false,
				Token:   "",
			}, nil
		}
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &pb.AuthenticateResponse{
		Success: true,
		Token:   token,
	}, nil
}
