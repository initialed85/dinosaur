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
	language     string
	mu           sync.Mutex
	uuid         uuid.UUID
	dead         bool
	code         string
	dockerRunCmd *exec.Cmd
	host         string
	port         int
	buildCmd     string
	runCmd       string
	sourcePath   string
	heartbeat    time.Time
}

func NewSession(language string) *Session {
	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err) // TODO shouldn't fail unless things are dire
	}

	s := Session{
		language: language,
		uuid:     sessionUUID,
		host:     fmt.Sprintf("session-%v-%v", language, sessionUUID.String()),
		port:     sessionPort,
		dead:     false,
	}

	return &s
}

func (s *Session) Language() string {
	return s.language
}

func (s *Session) UUID() uuid.UUID {
	return s.uuid
}

func (s *Session) Host() string {
	return s.host
}

func (s *Session) Port() int {
	return s.port
}

func (s *Session) Code() string {
	return s.code
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

func (s *Session) Dead() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.dead || time.Now().Sub(s.heartbeat) > sessionTimeout
}

func (s *Session) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []byte
	var err error

	supportedLanguage, ok := supportedLanguageByName[s.language]
	if !ok {
		return fmt.Errorf("unsupported language: %v", s.language)
	}

	s.buildCmd = supportedLanguage.BuildCmd
	s.runCmd = supportedLanguage.RunCmd
	s.sourcePath = path.Join(supportedLanguage.FolderPath, supportedLanguage.FileName)
	s.code = supportedLanguage.Code

	dockerRunCmd := fmt.Sprintf(
		`docker run --rm -t --cpus 0.1 --memory 0.1g --name %v --hostname %v --network dinosaur-internal --cap-add SYS_PTRACE -e GOTTY_PATH="%v" -e SESSION_UUID="%v" -e BUILD_CMD="%v" -e RUN_CMD="%v" dinosaur-session`,
		s.host,
		s.host,
		fmt.Sprintf("/proxy_session/%v/", s.uuid.String()),
		s.uuid.String(),
		s.buildCmd,
		s.runCmd,
	)

	s.dockerRunCmd = exec.Command(
		"bash",
		"-c",
		dockerRunCmd,
	)

	go func() {
		log.Printf("starting dockerRunCmd=%v", s.dockerRunCmd)
		out, err = s.dockerRunCmd.CombinedOutput()
		err = fmt.Errorf(string(out))
		log.Printf("stopped dockerRunCmd=%v", s.dockerRunCmd)
		s.dead = true
	}()

	runtime.Gosched()

	time.Sleep(time.Second * 2) // TODO: need something smarter than a magic sleep

	if err != nil {
		s.Close()
		return err
	}

	s.heartbeat = time.Now()

	return nil
}

func (s *Session) PushToSession(data string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.CreateTemp("", s.host)
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
			s.host,
			destPath,
		),
	)

	err = dockerCpCmd.Run()
	if err != nil {
		return err
	}

	s.code = data

	return nil
}

func (s *Session) Heartbeat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.heartbeat = time.Now()
}

func (s *Session) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	defer func() {
		s.dead = true
	}()

	if s.dockerRunCmd != nil && s.dockerRunCmd.Process != nil {
		_ = s.dockerRunCmd.Process.Kill()
	}

	dockerStopCmd := exec.Command("bash", "-c", fmt.Sprintf("docker kill %v", s.host))
	_ = dockerStopCmd.Run()

	dockerRmCmd := exec.Command("bash", "-c", fmt.Sprintf("docker rm -f %v", s.host))
	_ = dockerRmCmd.Run()
}
