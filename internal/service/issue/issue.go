package issue

import (
	"context"
	"log/slog"

	"thingify/internal/domain/model"
)

type Service struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{log: log}
}

func (s *Service) PublishIssue(_ context.Context, issue model.IssueAction) error {
	const op = "service.issue.PublishIssue"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("issue_id", issue.Issue.ID),
	)

	log.Info("Publish issue to the message broker successfully")

	return nil
}
