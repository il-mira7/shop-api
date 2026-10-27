package di

import (
	"shopAPI/internal/shop/authclient"
	"shopAPI/internal/shop/config"
	"shopAPI/internal/shop/handler"
	"shopAPI/internal/shop/repository/postgres"
	"shopAPI/internal/shop/service"

	"go.uber.org/fx"
)

func newUserClient(c *authclient.Client) handler.UserClient {
	return c
}

var PostgresModule = fx.Module("postgres",
	fx.Provide(
		config.NewConfig,
		postgres.NewPostgresPool,

		fx.Annotate(postgres.NewSupplierRepository, fx.As(new(service.SupplierRepository))),
		fx.Annotate(postgres.NewProductRepository, fx.As(new(service.ProductRepository))),
		fx.Annotate(postgres.NewImageRepository, fx.As(new(service.ImageRepository))),
		fx.Annotate(postgres.NewClientRepository, fx.As(new(service.ClientRepository))),
	),
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		fx.Annotate(service.NewSupplierService, fx.As(new(handler.SupplierService))),
		fx.Annotate(service.NewProductService, fx.As(new(handler.ProductService))),
		fx.Annotate(service.NewImageService, fx.As(new(handler.ImageService))),
		fx.Annotate(service.NewClientService, fx.As(new(handler.ClientService))),
	),
)

var HandlerModule = fx.Module("handler",
	fx.Provide(
		handler.NewSupplierHandler,
		handler.NewProductHandler,
		handler.NewImageHandler,
		handler.NewClientHandler,
		handler.NewUserHandler,
		authclient.NewClient,
		fx.Annotate(newUserClient, fx.As(new(handler.UserClient))),
	),
)
