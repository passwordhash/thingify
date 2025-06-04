package monitor

import (
	"context"
	"log/slog"
	"thingify/internal/domain/model"
	"thingify/internal/github"
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

var IssuesDB []model.Issue // TODO: убрать глобальную переменную

// ShortPollingNewIssues возвращает новые задачи пользователя, используя короткое опрос
func (m *Service) ShortPollingNewIssues(
	ctx context.Context,
	userToken string,
	pollingInterval time.Duration,
) ([]model.Issue, error) {
	const op = "monitor.ShortPollingNewIssues"

	log := m.log.With("op", op)

	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	log.InfoContext(ctx, "starting short polling for new issues", "pollingInterval", pollingInterval)

	// TODO: context with timeout

	pulledIssues := make(chan []model.Issue) // TODO: может лучше ссылка на Issue?
	newIssues := make(chan model.Issue, 5)

	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				// TODO: не забыть про обработку ошибок
				i, _ := m.issuesProvider.UserIssues(ctx, userToken)

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

				// for _, issue := range issues {
				// ...
				// }
				for _, newIssue := range issues {
					// if lastIssue := lastIssueCreatedAt(IssuesDB); lastIssue != nil {
					// ...
					newIssues <- newIssue
				}

			case <-ctx.Done():
				log.InfoContext(ctx, "stopping processing pulled issues")
				close(newIssues)
				return
			}
		}
	}()

	for newIssue := range newIssues {
		log.InfoContext(ctx, "new issue found", "issueID", newIssue.ID, "title", newIssue.Title, "createdAt", newIssue.CreatedAt)
	}

	// for issues := range issueSlices {
	// 	// if len(issue) > 0 {
	// 	// 	log.Info("new issues found", "count", len(issue))
	// 	// return issue, nil
	// 	// }
	// 	log.Info("new issues found", "count", len(issues))
	// 	for _, issue := range issues {
	// 		log.Info("issue found", "issue", issue.ID, "title", issue.Title, "createdAt", issue.CreatedAt)
	// 	}
	// }

	return nil, nil // TODO: return issues
}

func lastIssueCreatedAt(issues []model.Issue) *model.Issue {
	if len(issues) == 0 {
		return nil
	}

	lastIssue := issues[0]
	for _, issue := range issues {
		t1, err := time.Parse("2025-06-02T16:39:26Z", issue.CreatedAt)
		if err != nil {
			panic(err)
		}

		t, err := time.Parse("2025-06-02T16:39:26Z", lastIssue.CreatedAt)
		if err != nil {
			panic(err)
		}

		// if issue.CreatedAt.After(lastIssue.CreatedAt) {
		if t1.After(t) {
			lastIssue = issue
		}
	}

	return &lastIssue
}
