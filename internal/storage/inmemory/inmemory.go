package inmemory

import (
	"context"
	"sync"
	"thingify/internal/domain/model"
)

type Storage struct {
	mu     sync.RWMutex
	issues map[string][]model.Issue
}

func New() *Storage {
	return &Storage{
		issues: make(map[string][]model.Issue),
	}
}

func (s *Storage) Save(ctx context.Context, login string, issue model.Issue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.issues[login] = append(s.issues[login], issue)

	return nil
}

func (s *Storage) Issues(ctx context.Context, login string) ([]model.Issue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIssues, exists := s.issues[login]
	if !exists {
		return nil, nil // or return an error if preferred
	}

	return userIssues, nil
}
