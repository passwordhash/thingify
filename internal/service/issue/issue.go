package issue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"thingify/internal/domain/model"
	repoerr "thingify/internal/storage/errors"
)

type issuePublisher interface {
	Publish(ctx context.Context, routingKey string, payload any) error
}

type userIDSaver interface {
	SaveUserID(ctx context.Context, installationID, userID string) error
}

type userIDProvider interface {
	GetInstallationIDByUserID(ctx context.Context, userID string) (string, error)
}

type Service struct {
	log            *slog.Logger
	publisher      issuePublisher
	userIDSaver    userIDSaver
	userIDProvider userIDProvider
}

func New(
	log *slog.Logger,
	publisher issuePublisher,
	userIDSaver userIDSaver,
	userIDProvider userIDProvider,
) *Service {
	return &Service{
		log:            log,
		publisher:      publisher,
		userIDSaver:    userIDSaver,
		userIDProvider: userIDProvider,
	}
}

func (s *Service) InstallNewUser(ctx context.Context, userID, installationID string) error {
	const op = "service.issue.InstallNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
		slog.String("installation_id", installationID),
	)

	err := s.userIDSaver.SaveUserID(ctx, installationID, userID)
	if err != nil {
		log.Error("Failed to save user ID", slog.Any("error", err))

		return fmt.Errorf("%s: %v", op, err)
	}

	log.Info("User ID saved successfully")

	return nil
}

func (s *Service) PublishIssue(ctx context.Context, issue model.IssueAction) error {
	const op = "service.issue.PublishIssue"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("issue_id", issue.Issue.ID),
	)

	installID, err := s.userIDProvider.GetInstallationIDByUserID(ctx, strconv.FormatInt(issue.Sender.ID, 10))
	if errors.Is(err, repoerr.ErrInstallationIDNotFound) {
		log.Warn("Installation ID not found for user ID", slog.String("user_id", strconv.FormatInt(issue.Sender.ID, 10)))

		// TODO: заменить на ошибку, которая будет обработана на уровне контроллера
		return fmt.Errorf("%s: %v", op, err)
	}
	if err != nil {
		log.Error("Failed to get installation ID by user ID", slog.Any("error", err))

		return fmt.Errorf("%s: %v", op, err)
	}

	err = s.publisher.Publish(ctx, installID, issue)
	if err != nil {
		log.Error("Failed to publish issue to the message broker", slog.Any("error", err))

		return err
	}

	log.Info("Publish issue to the message broker successfully")

	return nil
}
