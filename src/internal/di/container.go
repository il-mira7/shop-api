package di

import (
	"shopAPI/internal/config"
	"shopAPI/internal/handler"
	"shopAPI/internal/repository/postgres"
	"shopAPI/internal/service"

	"go.uber.org/fx"
)

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
	),
)
