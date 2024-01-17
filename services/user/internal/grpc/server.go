package grpcAuth

import (
	"context"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"userService/internal/domain/entity"
	"userService/internal/service"
	"userService/pkg/logger"
)

type serverAPI struct {
	nordew.UnimplementedUserServer
	service service.UserService

	logger logger.Logger
}

func Register(gRPCServer *grpc.Server, service service.UserService) {
	nordew.RegisterUserServer(gRPCServer, &serverAPI{service: service})
}

func (s *serverAPI) SignUp(ctx context.Context, reqInput *nordew.SignUpRequest) (*nordew.Empty, error) {
	s.logger.Info(reqInput.GetUsername(), reqInput.GetEmail(), reqInput.GetPassword())

	input := entity.SignUpInput{
		Username: reqInput.GetUsername(),
		Email:    reqInput.GetEmail(),
		Password: reqInput.GetPassword(),
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
