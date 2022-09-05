package sessions

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Session struct {
	mu            sync.Mutex
	language      string
	code          string
	uuid          uuid.UUID
	port          int
	dockerRunCmd  *exec.Cmd
	dead          bool
	buildCmd      string
	runCmd        string
	containerName string
	sourcePath    string
	heartbeat     time.Time
}

func NewSession(language string, code string) *Session {
	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err) // TODO shouldn't fail unless things are dire
	}

	s := Session{
		language: language,
		code:     code,
		uuid:     sessionUUID,
		port:     8080,
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
	return s.containerName
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
	f, err := os.CreateTemp("", s.containerName)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	destPath := path.Join("/srv", s.sourcePath)

	dockerCpCmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"docker cp %v %v:%v",
			f.Name(),
			s.containerName,
			destPath,
		),
	)

	err = dockerCpCmd.Run()
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

	supportedLanguage, ok := supportedLanguageByName[s.language]
	if !ok {
		return fmt.Errorf("unsupported language: %v", s.language)
	}

	s.buildCmd = supportedLanguage.BuildCmd
	s.runCmd = supportedLanguage.RunCmd
	s.containerName = fmt.Sprintf("session-%v", s.uuid.String())
	s.sourcePath = path.Join(supportedLanguage.FolderPath, supportedLanguage.FileName)

	cmd := fmt.Sprintf(
		`docker run --rm --cpus 0.5 --memory 0.5g --name %v --network dinosaur-internal -p %v:8080/tcp -e GOTTY_PATH="%v" -e BUILD_CMD="%v" -e RUN_CMD="%v" dinosaur-session`,
		s.containerName,
		fmt.Sprintf("%v", s.port),
		fmt.Sprintf("/proxy_session/%v/", s.uuid.String()),
		s.buildCmd,
		s.runCmd,
	)

	s.dockerRunCmd = exec.Command(
		"bash",
		"-c",
		cmd,
	)

	log.Printf("executing cmd=%v", s.dockerRunCmd)

	go func() {
		err = s.dockerRunCmd.Run()
		s.dead = true
	}()

	runtime.Gosched()

	time.Sleep(time.Millisecond * 1000) // TODO: wait for ready w/ smart check vs suspicious sleep

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

	if s.dockerRunCmd != nil && s.dockerRunCmd.Process != nil {
		_ = s.dockerRunCmd.Process.Kill()
	}

	dockerStopCmd := exec.Command("bash", "-c", fmt.Sprintf("docker kill %v", s.containerName))
	_ = dockerStopCmd.Run()

	dockerRmCmd := exec.Command("bash", "-c", fmt.Sprintf("docker rm -f %v", s.containerName))
	_ = dockerRmCmd.Run()
}
