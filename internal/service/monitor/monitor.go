package monitor

import (
	"context"
	"log/slog"
	"thingify/internal/domain/model"
	"time"
)

type CandidateIssuesProvider interface {
	UserIssues(ctx context.Context, userToken string, limit int) ([]model.Issue, error)
}

type IssuesSaver interface {
	Save(ctx context.Context, login string, issue model.Issue) error
}

type IssuesProvider interface {
	Issues(ctx context.Context, login string) ([]model.Issue, error)
}

type IssuesPublisher interface {
	Publish(ctx context.Context, login string, issue any) error
}

// RequestsConsumer представляет собой интерфейс для потребления запросов на новые задачи.
type RequestsConsumer interface {
	// Consume начинает прослушивание канала запросов и возвращает каналы
	// для получения данных и ошибок.
	Consume(ctx context.Context) (<-chan []byte, <-chan error)

	// Close закрывает канал потребителя и освобождает ресурсы.
	Close() error
}

type Service struct {
	log *slog.Logger

	candidateIssuesProvider CandidateIssuesProvider
	issuesSaver             IssuesSaver
	issuesProvider          IssuesProvider
	issuesPublisher         IssuesPublisher
	reqConsumer             RequestsConsumer
}

func New(
	log *slog.Logger,
	candidateIssuesProvider CandidateIssuesProvider,
	issuesSaver IssuesSaver,
	issuesProvider IssuesProvider,
	issuesPublisher IssuesPublisher,
	reqConsumer RequestsConsumer,
) *Service {
	return &Service{
		log:                     log,
		candidateIssuesProvider: candidateIssuesProvider,
		issuesSaver:             issuesSaver,
		issuesProvider:          issuesProvider,
		issuesPublisher:         issuesPublisher,
		reqConsumer:             reqConsumer,
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
				i, _ := m.candidateIssuesProvider.UserIssues(ctx, userToken, 30)

				pulledIssues <- i
			case <-ctx.Done():
				// TODO: скорее всегда нужно будет завершать толко получение новых задач
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
					// TODO: TEMP: логин — токен
					// TODO: обработка ошибок
					dbIssues, err := m.issuesProvider.Issues(ctx, userToken)
					if err != nil {
						log.ErrorContext(ctx, "error getting issues from DB", "err", err)
					}

					latestDBTimeIssue := latestTimeIssue(dbIssues)
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
		m.issuesSaver.Save(ctx, userToken, newIssue)
		m.issuesPublisher.Publish(ctx, "tmp-login", newIssue)
		log.InfoContext(ctx, "new issue found", "issueID", newIssue.ID, "title", newIssue.Title, "createdAt", newIssue.CreatedAt)
	}

	return nil, nil // TODO: return issues
}

func (m *Service) StartRequestListener(ctx context.Context) {
	const op = "monitor.StartRequestListener"

	log := m.log.With("op", op)

	dataCh, errCh := m.reqConsumer.Consume(ctx)
	go func() {
		defer m.reqConsumer.Close()
		for {
			select {
			case <-ctx.Done():
				log.InfoContext(ctx, "stopping request listener")
				return
			case msg, ok := <-dataCh:
				if !ok {
					m.log.InfoContext(ctx, "request listener channel closed")
					return
				}

				log.InfoContext(ctx, "received request", "message", string(msg))
			case errCh, ok := <-errCh:
				if !ok {
					log.InfoContext(ctx, "request listener error channel closed")
					return
				}

				log.ErrorContext(ctx, "error in request listener", "error", errCh)
				return
			}
		}
	}()
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
