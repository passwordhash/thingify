package webhook

import (
	"thingify/internal/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) RegisterRoutes(g fiber.Router) {
	webhookG := g.Group("/webhook", middleware.ValidateHubSignature(h.secret))
	{
		issueG := webhookG.Group("/issue")
		{
			issueG.Post("", h.retrieveIssue)
		}
	}
}
