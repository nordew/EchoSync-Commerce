package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
	"userService/pkg/logger"
)

type jwtAuthenticator struct {
	signKey string
	logger  logger.Logger
}

func NewAuth(signKey string, logger logger.Logger) Authenticator {
	return &jwtAuthenticator{
		logger:  logger,
		signKey: signKey,
	}
}

type TokenClaims struct {
	UserId string `json:"sub"`
	jwt.RegisteredClaims
}

func (s *jwtAuthenticator) GenerateTokens(options *GenerateTokenClaimsOptions) (string, string, error) {
	mySigningKey := []byte(s.signKey)

	claims := TokenClaims{
		UserId: options.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "upload-api",
			Subject:   "client",
			ID:        uuid.NewString(),
			Audience:  []string{"upload"},
		},
	}

	refreshToken, err := s.GenerateRefreshToken(options.UserId)
	if err != nil {
		s.logger.Error("failed to generate refresh token", err)
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(mySigningKey)
	if err != nil {
		s.logger.Error("failed to generate access token", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *jwtAuthenticator) GenerateRefreshToken(id string) (string, error) {
	mySigningKey := []byte(s.signKey)

	claims := TokenClaims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "EchoSync",
			Subject:   "client",
			ID:        uuid.NewString(),
			Audience:  []string{"upload"},
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedRefreshToken, err := refreshToken.SignedString(mySigningKey)
	if err != nil {
		s.logger.Error("failed to generate refresh token", err)
		return "", err
	}

	return signedRefreshToken, nil
}

func (s *jwtAuthenticator) ParseToken(accessToken string) (*ParseTokenClaimsOutput, error) {
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error("unexpected signing method", nil)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		s.logger.Error("failed to parse jwt token", err)
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	if !token.Valid {
		s.logger.Error("token is not valid", nil)
		return nil, fmt.Errorf("token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role == nil {
		return nil, fmt.Errorf("token is not valid")
	}
	sub := claims["sub"]
	if sub == nil {
		return nil, fmt.Errorf("token is not valid")
	}

	return &ParseTokenClaimsOutput{Sub: fmt.Sprint(sub)}, nil
}
