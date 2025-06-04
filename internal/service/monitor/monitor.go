package monitor

import (
	"context"
	"log/slog"
	"thingify/internal/domain/model"
	"thingify/internal/github"
	"time"
)

type IssuesProvider interface {
	UserIssues(ctx context.Context, userToken string, limit int) ([]model.Issue, error)
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

var IssuesDB []model.Issue // TODO: убрать глобальную переменную

// ShortPollingNewIssues возвращает новые задачи пользователя, используя короткое опрос
func (m *Service) ShortPollingNewIssues(
	ctx context.Context,
	userToken string,
	pollingInterval time.Duration,
) ([]model.Issue, error) {
	const op = "monitor.ShortPollingNewIssues"

	log := m.log.With("op", op)

	log.InfoContext(ctx, "starting short polling for new issues", "pollingInterval", pollingInterval)

	// TODO: context with timeout
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	pulledIssues := make(chan []model.Issue) // TODO: может лучше ссылка на Issue?
	newIssues := make(chan model.Issue, 5)

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				const limit = 30 // TODO: сделать параметром
				// TODO: не забыть про обработку ошибок
				i, _ := m.issuesProvider.UserIssues(ctx, userToken, 30)

				pulledIssues <- i
			case <-ctx.Done():
				log.InfoContext(ctx, "stopping short polling for new issues")
				close(pulledIssues)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case issues := <-pulledIssues:
				if len(issues) == 0 {
					continue
				}

				for _, candidateIssue := range issues {
					latestDBTimeIssue := latestTimeIssue(IssuesDB)
					if latestDBTimeIssue == nil {
						newIssues <- candidateIssue
						continue
					}

					if candidateIssue.CreatedAt.After(*latestDBTimeIssue) {
						newIssues <- candidateIssue
						continue
					}
				}

			case <-ctx.Done():
				log.InfoContext(ctx, "stopping processing pulled issues")
				close(newIssues)
				return
			}
		}
	}()

	for newIssue := range newIssues {
		IssuesDB = append(IssuesDB, newIssue) // TODO: mock db запись
		log.InfoContext(ctx, "new issue found", "issueID", newIssue.ID, "title", newIssue.Title, "createdAt", newIssue.CreatedAt)
	}

	return nil, nil // TODO: return issues
}

// latestTimeIssue находит самую последнюю задачу по времени создания.
// Если список задач пуст, возвращает nil.
func latestTimeIssue(issues []model.Issue) *time.Time {
	if len(issues) == 0 {
		return nil
	}

	latest := &issues[0].CreatedAt
	for i := 1; i < len(issues); i++ {
		if issues[i].CreatedAt.After(*latest) {
			latest = &issues[i].CreatedAt
		}
	}
	return latest
}
