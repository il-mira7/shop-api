package grpc

import (
	"context"
	"errors"

	"shopAPI/internal/authserver/domain"
	"shopAPI/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	CreateUser(ctx context.Context, u domain.User, password string) (string, error)
	Authenticate(ctx context.Context, email string, password string) (string, error)
	ChangePassword(ctx context.Context, email string, passwordLast string, passwordNew string) (bool, error)
	ResetPassword(ctx context.Context, email string) (bool, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type authServer struct {
	pb.UnimplementedAuthServiceServer
	auth AuthService
}

func NewAuthServer(s AuthService) *authServer {
	return &authServer{auth: s}
}

func (s *authServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := domain.User{
		Email:       req.Email,
		Name:        req.Name,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
	}
	token, err := s.auth.CreateUser(ctx, user, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.CreateUserResponse{Token: token}, nil
}

func (s *authServer) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	token, err := s.auth.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.AuthenticateResponse{Token: token}, nil
}

func (s *authServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	res, err := s.auth.ChangePassword(ctx, req.Email, req.PasswordLast, req.PasswordNew)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.ChangePasswordResponse{Success: res}, nil
}

func (s *authServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	res, err := s.auth.ResetPassword(ctx, req.Email)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.ResetPasswordResponse{Success: res}, nil
}

func (s *authServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	id, err := s.auth.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.ValidateTokenResponse{Valid: true, UserId: id}, nil
}

func mapError(err error) error {
	var notFound *domain.NotFoundError
	var emptyPayload *domain.EmptyPayloadError
	var alreadyExists *domain.AlreadyExistsError
	switch {
	case errors.As(err, &notFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.As(err, &emptyPayload):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.As(err, &alreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return status.Error(codes.Internal, "internal error")
}
