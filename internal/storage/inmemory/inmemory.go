package inmemory

import (
	"context"
	"sync"
	repoerr "thingify/internal/storage/errors"
)

type Storage struct {
	mu sync.RWMutex

	// "userID": "installationID"
	users map[string]string
}

func New() *Storage {
	return &Storage{
		users: make(map[string]string),
	}
}

func (s *Storage) SaveUserID(_ context.Context, installationID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[userID] = installationID

	return nil
}

func (s *Storage) GetInstallationIDByUserID(_ context.Context, userID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.users[userID]
	if !ok {
		return "", repoerr.ErrInstallationIDNotFound
	}

	return id, nil
}
