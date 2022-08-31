package sessions

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Manager struct {
	mu            sync.Mutex
	sessionByUUID map[uuid.UUID]*Session
	ticker        *time.Ticker
}

func NewManager() *Manager {
	m := Manager{
		mu:            sync.Mutex{},
		sessionByUUID: make(map[uuid.UUID]*Session),
	}

	return &m
}

func (m *Manager) Open() error {
	m.ticker = time.NewTicker(time.Second)

	go func() {
		for {
			_ = <-m.ticker.C

			m.mu.Lock()
			dead := m.ticker == nil
			m.mu.Unlock()

			if dead {
				return
			}

			m.PruneSessions()
		}
	}()

	return nil
}

func (m *Manager) CreateSession(language string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	supportedLanguage, ok := supportedLanguageByName[language]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %v", language)
	}

	s := NewSession(language, supportedLanguage.Code)

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

func (m *Manager) PruneSessions() {
	m.mu.Lock()
	toDelete := make([]uuid.UUID, 0)
	for sessionUUID, session := range m.sessionByUUID {
		if !session.Dead() {
			continue
		}

		toDelete = append(toDelete, sessionUUID)
	}
	m.mu.Unlock()

	for _, sessionUUID := range toDelete {
		_ = m.DestroySession(sessionUUID)
	}
}

func (m *Manager) Close() {
	m.mu.Lock()
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}
	toDelete := make([]uuid.UUID, 0)
	for sessionUUID, _ := range m.sessionByUUID {
		toDelete = append(toDelete, sessionUUID)
	}
	m.mu.Unlock()

	for _, sessionUUID := range toDelete {
		_ = m.DestroySession(sessionUUID)
	}
}
