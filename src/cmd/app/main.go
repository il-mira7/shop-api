package main

import (
	"log"
	_ "shopAPI/docs"
	"shopAPI/internal/di"
	"shopAPI/internal/handler"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

//	@title			Swagger Shop API (BE2_APIFromSpec)
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/api/v1

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	fx.New(
		di.PostgresModule,
		di.ServiceModule,
		di.HandlerModule,
		fx.Invoke(handler.NewHTTPServer),
	).Run()
}
