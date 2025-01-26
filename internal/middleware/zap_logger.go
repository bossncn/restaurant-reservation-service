package middleware

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strings"
	"time"
)

func ZapLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeStarted := time.Now()

			path := c.Request().URL.Path
			if path != "/" && !strings.Contains(path, "swagger") {
				// Log incoming request
				logger.Info("Incoming request",
					zap.String("method", c.Request().Method),
					zap.String("path", path),
					zap.String("remote_addr", c.Request().RemoteAddr),
				)
			}

			// Process the request
			err := next(c)

			status := c.Response().Status
			httpErr := new(echo.HTTPError)
			if errors.As(err, &httpErr) {
				status = httpErr.Code
			}

			if path != "/" && !strings.Contains(path, "swagger") {
				fields := []zap.Field{
					zap.String("method", c.Request().Method),
					zap.String("path", c.Request().URL.Path),
					zap.String("query", c.Request().URL.RawQuery),
					zap.Int("status", status),
					zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(time.Since(timeStarted).Microseconds())/1000)),
				}

				logger.Info("Request completed", fields...)
			}

			return err
		}
	}
}
