package http

import (
	"github.com/bossncn/restaurant-reservation-service/config"
	_ "github.com/bossncn/restaurant-reservation-service/docs"
	"github.com/bossncn/restaurant-reservation-service/internal/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"net/http"
)

type Repository struct {
}

type Middleware struct {
	Logger echo.MiddlewareFunc
}

type Handler struct {
}

type Service struct {
}

func InitRepository(cfg *config.Config) *Repository {
	return &Repository{}
}

func InitMiddleware(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		Logger: middleware.ZapLoggerMiddleware(logger),
	}
}

func InitHandler(cfg *config.Config) *Handler {
	return &Handler{}
}

func InitService(cfg *config.Config) *Service {
	return &Service{}
}

type ServerHttp struct {
	app *echo.Echo
}

func NewHTTPServer(cfg *config.Config, middleware *Middleware, repository *Repository, handler *Handler, service *Service) *ServerHttp {
	e := echo.New()
	e.Use(middleware.Logger)

	// Healthcheck
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Healthy") })

	// Swagger
	if cfg.AppEnv == "development" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return &ServerHttp{
		app: e,
	}
}

func (s *ServerHttp) Start() {
	port := ":8080"
	_ = s.app.Start(port)
}
