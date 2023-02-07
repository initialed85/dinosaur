package sessions

import (
	"fmt"
	"github.com/google/uuid"
	"sort"
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
	// TODO: I've disabled the build within the lifecycle of the Manager for now- we build it during deployment which is
	//   probably good enough; sure if you do a "docker image prune" then you'll be goosed, but meh
	//dockerBuildCmd := exec.Command(
	//	"bash",
	//	"-c",
	//	"docker build -t dinosaur-session -f ./docker/session/Dockerfile ./docker/session/",
	//)
	//
	//output, err := dockerBuildCmd.CombinedOutput()
	//if err != nil {
	//	log.Printf("STDOUT / STDERR: %v", string(output))
	//	return err
	//}

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

func (m *Manager) GetSupportedLanguages() []SupportedLanguage {
	supportedLanguages := make([]SupportedLanguage, 0)
	for _, supportedLanguage := range supportedLanguageByName {
		supportedLanguages = append(
			supportedLanguages,
			SupportedLanguage{
				Name:         supportedLanguage.Name,
				FriendlyName: supportedLanguage.FriendlyName,
			},
		)
	}

	sort.Slice(
		supportedLanguages,
		func(i, j int) bool {
			return supportedLanguages[i].FriendlyName < supportedLanguages[j].FriendlyName
		},
	)

	return supportedLanguages
}

func (m *Manager) CreateSession(language string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
