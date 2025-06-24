package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	RequestIDKey       = "request_id"
	RequestIDHeaderKey = "X-Request-ID"
)

func Logging(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestID := c.Get(RequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Context().SetUserValue(RequestIDKey, requestID)
		//c.Locals(RequestIDKey, requestID)

		err := c.Next()

		status := c.Response().StatusCode()

		//logger.WithRequest(requestID)
		fields := []slog.Attr{
			slog.String("method", c.Method()),
			slog.Int("status", status),
			slog.String("path", c.Path()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		}

		switch {
		case err != nil:
			fields = append(fields, slog.Any("error", err))
			logger.LogAttrs(c.Context(), slog.LevelError, "Request failed with handler error", fields...)
		case status >= 500:
			fields = append(fields, slog.Any("error", err))
			logger.LogAttrs(c.Context(), slog.LevelError, "Request failed with server error", fields...)
		case status >= 400:
			fields = append(fields, slog.Any("error", err))
			logger.LogAttrs(c.Context(), slog.LevelWarn, "Request completed with client error", fields...)
		default:
			logger.LogAttrs(c.Context(), slog.LevelInfo, "Request completed", fields...)
		}

		return nil
	}
}
