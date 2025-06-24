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
	issueEvent = "issues"
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

func (h *handler) retrieveIssue(c *fiber.Ctx) error {
	event := c.Get(eventHeader)

	fmt.Println("event: ", event)
	if event != issueEvent {
		return response.BadRequest(c,
			fmt.Errorf("unsupported event type: %s", event),
			"Unsupported event type",
		)
	}

	var issue issueWebhookReq
	if err := c.BodyParser(&issue); err != nil {
		return response.BadRequest(c,
			fmt.Errorf("failed to parse request body: %w", err),
			"Failed to parse request body",
		)
	}

	// TODO: может быть вынести в middleware
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
