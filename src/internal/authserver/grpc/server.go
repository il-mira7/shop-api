package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"shopAPI/pb"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GRPCApp struct {
	srv *grpc.Server
}

func NewGRPCApp(s *authServer) *GRPCApp {
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(LoggingInterceptor))
	pb.RegisterAuthServiceServer(srv, s)
	reflection.Register(srv) // для grpcurl

	return &GRPCApp{srv: srv}
}

func NewGRPCServer(lc fx.Lifecycle, app *GRPCApp) {

	lc.Append(fx.Hook{
		OnStart: app.Connect,
		OnStop: func(_ context.Context) error {
			app.srv.GracefulStop()
			slog.Info("gRPC-сервер остановлен")
			return nil
		},
	})
}

func (app *GRPCApp) Connect(_ context.Context) error {
	port := os.Getenv("AUTH_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("listen :%s: %w", port, err)
	}
	slog.Info("gRPC-сервер запущен", "addr", lis.Addr())

	go func() {
		if err := app.srv.Serve(lis); err != nil {
			slog.Error("gRPC-сервер упал", "error", err)
		}
	}()
	return nil
}

// для красоты в терминале, как у гин
var (
	grpcGreen  = "\033[97;42m"
	grpcYellow = "\033[90;43m"
	grpcRed    = "\033[97;41m"
	grpcReset  = "\033[0m"
)

func codeColor(code codes.Code) string {
	switch code {
	case codes.OK:
		return grpcGreen
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition:
		return grpcYellow
	default:
		return grpcRed
	}
}

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	addr := "unknown"
	if p, ok := peer.FromContext(ctx); ok && p != nil {
		addr = p.Addr.String()
	}
	fmt.Fprintf(os.Stdout, "[GRPC] %s |%s %-15s %s| %13v | %22s | %s\n",
		start.Format("2006/01/02 - 15:04:05"),
		codeColor(status.Code(err)), status.Code(err), grpcReset,
		time.Since(start),
		addr,
		info.FullMethod,
	)
	if err != nil {
		fmt.Fprintf(os.Stdout, "       └─ %v\n", err)
	}

	return resp, err
}
