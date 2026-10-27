package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "shopAPI/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func NewHTTPServer(h *supplierHandler, i *imageHandler, p *productHandler, c *clientHandler, lifecycle fx.Lifecycle) *http.Server {
	r := gin.Default()
	setupMiddleware(r)

	server := setupRoutes(r, h, i, p, c)

	setupLifecycle(server, lifecycle)
	return server
}

func setupMiddleware(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
}

func setupRoutes(r *gin.Engine, s *supplierHandler, i *imageHandler, p *productHandler, c *clientHandler) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
	{
		suppliers := api.Group("/suppliers")
		{
			suppliers.POST("", s.Create)
			suppliers.GET("", s.GetAll)
			suppliers.GET("/:id", s.GetByID)
			suppliers.PATCH("/:id/address", s.Update)
			suppliers.DELETE("/:id", s.Delete)
		}

		images := api.Group("/images")
		{
			images.POST("/product/:product_id", i.Create)
			images.PUT("/:id", i.Update)
			images.DELETE("/:id", i.Delete)
			images.GET("/product/:product_id", i.GetByProductID)
			images.GET("/:id", i.GetByImageID)
		}

		products := api.Group("/products")
		{
			products.POST("", p.Create)
			products.PATCH("/:id/stock", p.DecreaseStock)
			products.GET("/:id", p.GetByID)
			products.GET("", p.GetAll)
			products.DELETE("/:id", p.Delete)
		}

		clients := api.Group("/clients")
		{
			clients.POST("", c.Create)
			clients.DELETE("/:id", c.Delete)
			clients.GET("/search", c.GetByNames)
			clients.GET("", c.GetAll)
			clients.PATCH("/:id/address", c.Update)
		}
	}

	return &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

}

func setupLifecycle(srv *http.Server, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			slog.Info("Сервер запущен ", "addr", srv.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("Критическая ошибка сервера:", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Завершение работы: закрываем соединения...")
			if err := srv.Shutdown(ctx); err != nil {
				slog.Error("Ошибка при остановке сервера:", "error", err)
				return fmt.Errorf("Ошибка при остановке сервера: %w", err)
			}
			slog.Info("Сервер успешно остановлен.")
			return nil
		},
	})
}
