package monitor

import (
	"context"
	"log/slog"
	"thingify/internal/github"
	"thingify/internal/model"
	"time"
)

type IssuesProvider interface {
	UserIssues(ctx context.Context, userToken string) ([]model.Issue, error)
}

type Service struct {
	log *slog.Logger

	issuesProvider IssuesProvider
}

func New(
	log *slog.Logger,
	ghClient *github.GHClient,
) *Service {
	return &Service{
		log:            log,
		issuesProvider: ghClient,
	}
}

// ShortPollingNewIssues возвращает новые задачи пользователя, используя короткое опрос
func (m *Service) ShortPollingNewIssues(
	ctx context.Context,
	userToken string,
	pollingInterval time.Duration,
) ([]model.Issue, error) {
	const op = "monitor.ShortPollingNewIssues"

	log := m.log.With("op", op)

	log.Info("starting short polling for new issues", "pollingInterval", pollingInterval)

	// TODO: context with timeout

	issueSlices := make(chan []model.Issue) // TODO: может лучше ссылка на Issue?

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				// TODO: не забыть про обработку ошибок
				i, _ := m.issuesProvider.UserIssues(ctx, userToken)

				issueSlices <- i
			}
		}
	}()

	for issues := range issueSlices {
		// if len(issue) > 0 {
		// 	log.Info("new issues found", "count", len(issue))
		// return issue, nil
		// }
		log.Info("new issues found", "count", len(issues))
		for _, issue := range issues {
			log.Info("issue found", "issue", issue.ID, "title", issue.Title, "createdAt", issue.CreatedAt)
		}
	}

	return nil, nil // TODO: return issues
}
