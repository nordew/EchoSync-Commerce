package grpcAuth

import (
	"context"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"userService/internal/domain/entity"
	"userService/internal/service"
)

type serverAPI struct {
	nordew.UnimplementedUserServer
	service service.UserService
}

func Register(gRPCServer *grpc.Server, service service.UserService) {
	nordew.RegisterUserServer(gRPCServer, &serverAPI{service: service})
}

func (s *serverAPI) SignUp(ctx context.Context, reqInput *nordew.SignUpRequest) (*nordew.Empty, error) {
	input := entity.SignUpInput{
		Username: reqInput.Username,
		Email:    reqInput.Email,
		Password: reqInput.Password,
	}

	if err := input.Validate(); err != nil {
		return nil, status.Error(400, "invalid input")
	}

	if err := s.service.SignUp(ctx, input); err != nil {
		return nil, status.Error(codes.Internal, "failed to sign up")
	}

	return &nordew.Empty{}, nil
}

func (s *serverAPI) SignIn(ctx context.Context, reqInput *nordew.SignInRequest) (*nordew.SignInResponse, error) {
	input := entity.SignInInput{
		Email:    reqInput.Email,
		Password: reqInput.Password,
	}

	if err := input.Validate(); err != nil {
		return nil, status.Error(400, "invalid input")
	}

	accessToken, refreshToken, err := s.service.SignIn(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign in")
	}

	resp := &nordew.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}