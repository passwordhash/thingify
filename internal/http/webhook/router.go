package webhook

import (
	"thingify/internal/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) RegisterRoutes(g fiber.Router) {
	webhookGroup := g.Group("/webhook", middleware.ValidateHubSignature(h.secret))

	webhookGroup.Post("/github", h.webhook)
}
