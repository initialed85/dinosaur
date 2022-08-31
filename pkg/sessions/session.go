package sessions

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Session struct {
	mu         sync.Mutex
	language   string
	code       string
	uuid       uuid.UUID
	port       int
	gottyCmd   *exec.Cmd
	dead       bool
	basePath   string
	folderPath string
	filePath   string
	buildCmd   string
	heartbeat  time.Time
}

func NewSession(language string, code string) *Session {
	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err) // TODO shouldn't fail unless things are dire
	}

	port, err := getFreePort()
	if err != nil {
		log.Fatal(err) // TODO shouldn't fail unless things are dire
	}

	s := Session{
		language: language,
		code:     code,
		uuid:     sessionUUID,
		port:     port,
		dead:     false,
	}

	return &s
}

func (s *Session) Language() string {
	return s.language
}

func (s *Session) Code() string {
	return s.code
}

func (s *Session) UUID() uuid.UUID {
	return s.uuid
}

func (s *Session) Host() string {
	return "localhost"
}

func (s *Session) Port() int {
	return s.port
}

func (s *Session) InternalURL() string {
	return fmt.Sprintf("http://%v:%v/proxy_session/%v/", s.Host(), s.Port(), s.UUID().String())
}

func (s *Session) GetProxyURL(externalURL *url.URL) string {
	return fmt.Sprintf(
		"http://%v:%v/%v",
		s.Host(),
		s.Port(),
		strings.TrimLeft(externalURL.Path, "/"),
	)
}

func (s *Session) PushToSession(data string) error {
	err := os.WriteFile(s.filePath, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Heartbeat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.heartbeat = time.Now()
}

func (s *Session) Dead() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.dead || time.Now().Sub(s.heartbeat) > time.Second*5
}

func (s *Session) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	s.basePath = filepath.Join(cwd, "tmp", s.uuid.String())

	supportedLanguage, ok := supportedLanguageByName[s.language]
	if !ok {
		return fmt.Errorf("unsupported language: %v", s.language)
	}

	s.folderPath = filepath.Join(s.basePath, supportedLanguage.FolderPath)
	s.filePath = filepath.Join(s.folderPath, supportedLanguage.FileName)
	s.buildCmd = supportedLanguage.BuildCmd

	err = os.MkdirAll(s.folderPath, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, []byte{}, 0644)
	if err != nil {
		return err
	}

	// TODO introduce the Docker layer somewhere around here
	cmd := fmt.Sprintf(
		`gotty --address 0.0.0.0 --port %v --path %v --ws-origin '.*' bash -c 'cd %v && find . -type f | entr -n -r -a -c -s "%v"'`,
		fmt.Sprintf("%v", s.port),
		fmt.Sprintf("/proxy_session/%v/", s.uuid.String()),
		s.basePath,
		s.buildCmd,
	)

	s.gottyCmd = exec.Command(
		"bash",
		"-c",
		cmd,
	)

	log.Printf("executing cmd=%v", s.gottyCmd)

	go func() {
		err = s.gottyCmd.Run()
		s.dead = true
	}()

	runtime.Gosched()

	time.Sleep(time.Millisecond * 100)

	if err != nil {
		s.Close()
		return err
	}

	s.heartbeat = time.Now()

	return nil
}

func (s *Session) Close() {
	defer func() {
		s.dead = true
	}()

	if s.gottyCmd != nil && s.gottyCmd.Process != nil {
		_ = s.gottyCmd.Process.Kill()
	}

	_ = os.RemoveAll(filepath.Join("tmp", s.uuid.String()))
}
