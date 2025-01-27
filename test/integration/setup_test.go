package integration_test

import (
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/event"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/http"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Setup() *echo.Echo {
	logger := zap.NewNop()
	e := echo.New()
	repo := http.InitRepository()
	eventProcessor, requestEvent := event.NewProcessor(repo.TableRepository, repo.ReservationRepository, logger)
	go eventProcessor.ProcessRequests()
	service := http.InitService(logger, repo, requestEvent)
	handlers := http.InitHandler(logger, service)
	handlers.TableHandler.RegisterRoutes(e.Group("/public"))
	handlers.ReservationHandler.RegisterRoutes(e.Group("/public"), e.Group("/secure"))

	return e
}
