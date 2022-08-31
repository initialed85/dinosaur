package sessions

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Manager struct {
	mu            sync.Mutex
	sessionByUUID map[uuid.UUID]*Session
}

func NewManager() *Manager {
	m := Manager{
		mu:            sync.Mutex{},
		sessionByUUID: make(map[uuid.UUID]*Session),
	}

	return &m
}

func (m *Manager) CreateSession(language string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO cleanup of sessions
	s := NewSession(language)

	err := s.Open()
	if err != nil {
		return nil, err
	}

	m.sessionByUUID[s.UUID()] = s

	return s, nil
}

func (m *Manager) GetSession(sessionUUID uuid.UUID) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessionByUUID[sessionUUID]
	if !ok {
		return nil, fmt.Errorf("no session for session UUID: %#+v", sessionUUID.String())
	}

	return s, nil
}

func (m *Manager) DestroySession(sessionUUID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessionByUUID[sessionUUID]
	if !ok {
		return fmt.Errorf("no session for session UUID: %#+v", sessionUUID.String())
	}

	s.Close()

	delete(m.sessionByUUID, sessionUUID)

	return nil
}

func (m *Manager) Close() {
	m.mu.Lock()
	toDelete := make([]uuid.UUID, 0)
	for sessionUUID, _ := range m.sessionByUUID {
		toDelete = append(toDelete, sessionUUID)
	}
	m.mu.Unlock()

	for _, sessionUUID := range toDelete {
		_ = m.DestroySession(sessionUUID)
	}
}
