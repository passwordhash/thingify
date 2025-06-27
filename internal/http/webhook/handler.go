package webhook

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"thingify/internal/domain/model"
	"thingify/internal/http/response"
	"thingify/internal/http/webhook/dto"

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
	case installationEvent:
		return h.handleInstallationEvent(c)
	default:
		return response.BadRequest(c,
			fmt.Errorf("unsupported event type: %s", event),
			"Unsupported event type",
		)
	}
}

func (h *handler) handleIssueEvent(c *fiber.Ctx) error {
	var issue dto.IssueWebhookReq
	if err := c.BodyParser(&issue); err != nil {
		return response.BadRequest(c,
			fmt.Errorf("failed to parse request body: %w", err),
			"Failed to parse request body",
		)
	}

	if issue.Action != dto.ActionOpened {
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

// 1.	Проверить подпись (безопасность) DONE
// 2.	Распарсить событие DONE
// 3.	Сохранить данные установки (installation ID + пользователь)
// 4.	(опционально) Получить installation access token
// 5.	(опционально) Получить список установленных репозиториев
// 6.	Привязать установку к своему пользователю (если надо)
// 7.	Подтвердить полученное событие (200 OK)

func (h *handler) handleInstallationEvent(c *fiber.Ctx) error {
	var inst dto.InstallationWebhookReq
	if err := c.BodyParser(&inst); err != nil {
		return response.BadRequest(c,
			fmt.Errorf("failed to parse request body: %w", err),
			"Failed to parse request body",
		)
	}

	log.Println("hello")

	// TODO: add delete handling
	if inst.Action != dto.ActionCreated {
		return response.BadRequest(c,
			fmt.Errorf("unsupported action: %s", inst.Action),
			"Unsupported action",
		)
	}

	fmt.Printf("%+v", inst)

	return nil
}
