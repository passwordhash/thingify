package issue

import (
	"context"
	"log/slog"

	"thingify/internal/domain/model"
)

type issuePublisher interface {
	Publish(ctx context.Context, routingKey string, payload any) error
}

type Service struct {
	log       *slog.Logger
	publisher issuePublisher
}

func New(log *slog.Logger, publisher issuePublisher) *Service {
	return &Service{
		log:       log,
		publisher: publisher,
	}
}

func (s *Service) PublishIssue(ctx context.Context, issue model.IssueAction) error {
	const op = "service.issue.PublishIssue"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("issue_id", issue.Issue.ID),
	)

	// TODO: сделать routingKey - installation id
	installID := "73271876"
	err := s.publisher.Publish(ctx, installID, issue)
	if err != nil {
		log.Error("Failed to publish issue to the message broker", slog.Any("error", err))

		return err
	}

	log.Info("Publish issue to the message broker successfully")

	return nil
}
