package http

import (
	"github.com/bossncn/restaurant-reservation-service/config"
	_ "github.com/bossncn/restaurant-reservation-service/docs"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/memory"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"github.com/bossncn/restaurant-reservation-service/internal/core/repository"
	"github.com/bossncn/restaurant-reservation-service/internal/core/service"
	"github.com/bossncn/restaurant-reservation-service/internal/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"net/http"
)

type Repository struct {
	TableRepository       repository.TableRepository
	ReservationRepository repository.ReservationRepository
}

type Middleware struct {
	Logger echo.MiddlewareFunc
}

type Handler struct {
	TableHandler       *TableHandler
	ReservationHandler *ReservationHandler
}

type Service struct {
	TableService       *service.TableService
	ReservationService *service.ReservationService
}

func InitRepository() *Repository {
	return &Repository{
		TableRepository:       memory.NewTableRepository(),
		ReservationRepository: memory.NewReservationRepository(),
	}
}

func InitMiddleware(logger *zap.Logger) *Middleware {
	return &Middleware{
		Logger: middleware.ZapLoggerMiddleware(logger),
	}
}

func InitHandler(logger *zap.Logger, services *Service) *Handler {
	return &Handler{
		TableHandler:       NewTableHandler(logger, services),
		ReservationHandler: NewReservationHandler(logger, services),
	}
}

func InitService(logger *zap.Logger, repo *Repository, eventRequest *chan model.EventRequest) *Service {
	return &Service{
		TableService:       service.NewTableService(repo.TableRepository, logger, eventRequest),
		ReservationService: service.NewReservationService(repo.ReservationRepository, logger, eventRequest),
	}
}

type ServerHttp struct {
	app *echo.Echo
}

func NewHTTPServer(cfg *config.Config, middleware *Middleware, handler *Handler) *ServerHttp {
	e := echo.New()
	e.Use(middleware.Logger)

	// Healthcheck
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Healthy") })

	// Swagger
	if cfg.AppEnv == "development" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	publicRoute := e.Group("/public")
	secureRoute := e.Group("/secure")
	handler.TableHandler.RegisterRoutes(publicRoute)
	handler.ReservationHandler.RegisterRoutes(publicRoute, secureRoute)

	return &ServerHttp{
		app: e,
	}
}

func (s *ServerHttp) Start() {
	port := ":8080"
	_ = s.app.Start(port)
}
