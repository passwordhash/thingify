package webhook

import (
	"context"
	"fmt"
	"net/http"

	"thingify/internal/domain/model"
	"thingify/internal/http/response"

	"github.com/gofiber/fiber/v2"
)

const (
	eventHeader = "X-GitHub-Event"
)

const (
	issueEvent        = "issues"
	installationEvent = "installation"
)

type issuePublisher interface {
	PublishIssue(ctx context.Context, issue model.IssueAction) error
}

type handler struct {
	issuePublisher issuePublisher
	secret         string
}

func NewHandler(
	publisher issuePublisher,
	secret string,
) *handler {
	return &handler{
		issuePublisher: publisher,
		secret:         secret,
	}
}

func (h *handler) webhook(c *fiber.Ctx) error {
	event := c.Get(eventHeader)

	switch event {
	case issueEvent:
		return h.handleIssueEvent(c)
	default:
		return response.BadRequest(c,
			fmt.Errorf("unsupported event type: %s", event),
			"Unsupported event type",
		)
	}
}

func (h *handler) handleIssueEvent(c *fiber.Ctx) error {
	var issue issueWebhookReq
	if err := c.BodyParser(&issue); err != nil {
		return response.BadRequest(c,
			fmt.Errorf("failed to parse request body: %w", err),
			"Failed to parse request body",
		)
	}

	if issue.Action != "opened" {
		return response.BadRequest(c,
			fmt.Errorf("unsupported action: %s", issue.Action),
			"Unsupported action",
		)
	}

	domain, err := issue.ToDomain()
	if err != nil {
		return response.BadRequest(c,
			fmt.Errorf("failed to convert issue to domain: %w", err),
			"Failed to convert issue to domain",
		)
	}

	err = h.issuePublisher.PublishIssue(c.UserContext(), domain)
	if err != nil {
		fmt.Println(3)
		return response.Internal(c, err, "Failed to publish issue")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"issue": domain,
	})
}
