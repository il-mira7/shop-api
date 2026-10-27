package authclient

import (
	"context"
	"fmt"

	"shopAPI/internal/shop/config"
	"shopAPI/internal/shop/domain"
	"shopAPI/pb"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	stub pb.AuthServiceClient
}

func NewClient(cfg *config.Config, lc fx.Lifecycle) (*Client, error) {
	conn, err := grpc.NewClient(cfg.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return conn.Close()
		},
	})

	return &Client{stub: pb.NewAuthServiceClient(conn)}, nil
}

func (c *Client) Register(ctx context.Context, email, name, surname, phone, password string) (string, error) {
	resp, err := c.stub.CreateUser(ctx, &pb.CreateUserRequest{
		Email: email, Name: name, Surname: surname, PhoneNumber: phone, Password: password,
	})
	return resp.GetToken(), mapGRPCError(err)
}

func (c *Client) Authenticate(ctx context.Context, email, password string) (string, error) {
	resp, err := c.stub.Authenticate(ctx, &pb.AuthenticateRequest{
		Email: email, Password: password,
	})
	return resp.GetToken(), mapGRPCError(err)
}

func (c *Client) ResetPassword(ctx context.Context, email string) (bool, error) {
	resp, err := c.stub.ResetPassword(ctx, &pb.ResetPasswordRequest{
		Email: email,
	})
	return resp.GetSuccess(), mapGRPCError(err)
}

func (c *Client) ValidateToken(ctx context.Context, token string) (string, error) {
	resp, err := c.stub.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		return "", mapGRPCError(err)
	}
	return resp.GetUserId(), nil
}

func mapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return &domain.NotFoundError{Msg: st.Message()}
	case codes.InvalidArgument:
		return &domain.ValidationError{Message: st.Message()}
	case codes.AlreadyExists:
		return &domain.AlreadyExistsError{Msg: st.Message()}
	case codes.Unauthenticated:
		return domain.ErrInvalidCredentials
	default:
		return fmt.Errorf("%w: %s", domain.ErrAuthServiceUnavailable, st.Message())
	}
}
