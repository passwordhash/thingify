package webhook

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	secret string
}

func NewHandler(secret string) *handler {
	return &handler{
		secret: secret,
	}
}

func (h *handler) retrieveIssue(c *fiber.Ctx) error {
	body := c.Body()

	out, err := os.Create("out.json")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create out file",
		})
	}
	defer out.Close()
	_, err = out.Write(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to write output",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"event": "X-GitHub-Event",
	})
}
