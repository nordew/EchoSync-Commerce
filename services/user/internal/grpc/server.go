package grpcAuth

import (
	"context"
	"fmt"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"unicode/utf8"
	"userService/internal/domain/entity"
	"userService/internal/service"
	"userService/pkg/logger"
)

type serverAPI struct {
	nordew.UnimplementedUserServer
	service service.UserService

	logger logger.Logger
}

func Register(gRPCServer *grpc.Server, service service.UserService, logger logger.Logger) {
	nordew.RegisterUserServer(gRPCServer, &serverAPI{service: service, logger: logger})
}

func (s *serverAPI) SignUp(ctx context.Context, reqInput *nordew.SignUpRequest) (*nordew.Empty, error) {
	s.logger.Info("converting input")
	input := entity.SignUpInput{
		Username: reqInput.GetUsername(),
		Email:    reqInput.GetEmail(),
		Password: reqInput.GetPassword(),
	}

	if err := validateSignUpInput(input.Username, input.Email, input.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s.logger.Info("signing up")
	if err := s.service.SignUp(ctx, &input); err != nil {
		return nil, status.Error(codes.Internal, "failed to sign up")
	}

	return &nordew.Empty{}, nil
}

func (s *serverAPI) SignIn(ctx context.Context, reqInput *nordew.SignInRequest) (*nordew.SignInResponse, error) {
	input := entity.SignInInput{
		Email:    reqInput.Email,
		Password: reqInput.Password,
	}

	if !isValidEmail(input.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	accessToken, refreshToken, err := s.service.SignIn(ctx, &input)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to sign in")
	}

	resp := &nordew.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func (s *serverAPI) RefreshToken(ctx context.Context, reqInput *nordew.RefreshRequest) (*nordew.RefreshResponse, error) {
	accessToken, refreshToken, err := s.service.RefreshTokens(ctx, reqInput.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	resp := &nordew.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return resp, nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func validateSignUpInput(username, email, password string) error {
	if utf8.RuneCountInString(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if utf8.RuneCountInString(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if !isValidEmail(email) {
		return fmt.Errorf("invalid email address")
	}

	return nil
}
