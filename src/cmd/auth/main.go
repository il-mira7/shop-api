package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"shopAPI/internal/authserver/config"
	"shopAPI/internal/authserver/grpc"
	"shopAPI/internal/authserver/repository/postgres"
	"shopAPI/internal/authserver/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	fx.New(
		fx.Provide(
			grpc.NewGRPCApp,
			config.NewConfig,
			postgres.NewPostgresPool,
			service.NewJwtProvider,
			grpc.NewAuthServer,
			//postgres.NewUserRepository,
			//service.NewAuthService,
			fx.Annotate(postgres.NewUserRepository, fx.As(new(service.UserRepository))), // репозиторий → интерфейс сервиса
			fx.Annotate(service.NewAuthService, fx.As(new(grpc.AuthService))),           // сервис → интерфейс хендлера

		),
		fx.Invoke(grpc.NewGRPCServer),
	).Run()
}
