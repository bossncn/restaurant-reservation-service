package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func ZapLoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeStarted := time.Now()

			// Log incoming request
			logger.Info("Incoming request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("remote_addr", c.Request().RemoteAddr),
			)

			// Process the request
			err := next(c)

			status := c.Response().Status
			httpErr := new(echo.HTTPError)
			if errors.As(err, &httpErr) {
				status = httpErr.Code
			}

			fields := []zap.Field{
				zap.Int64("elapsed", int64(time.Since(timeStarted)/time.Millisecond)),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("query", c.Request().URL.RawQuery),
				zap.Int("status", status),
			}

			logger.Info("Request completed", fields...)

			return err
		}
	}
}
