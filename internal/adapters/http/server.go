package http

import (
	"github.com/bossncn/go-boilerplate/config"
	"github.com/bossncn/go-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
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

func NewHTTPServer(middleware *Middleware, repository *Repository, handler *Handler, service *Service) *ServerHttp {
	e := echo.New()
	e.Use(middleware.Logger)
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Healthy") })

	return &ServerHttp{
		app: e,
	}
}

func (s *ServerHttp) Start() {
	port := ":8080"

	_ = s.app.Start(port)
}
