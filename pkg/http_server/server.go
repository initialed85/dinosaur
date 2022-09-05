package http_server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/initialed85/dinosaur/pkg/sessions"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

type Server struct {
	serveMux       *http.ServeMux
	server         *http.Server
	sessionManager *sessions.Manager
}

func New(
	port int,
	sessionManager *sessions.Manager,
) *Server {
	serveMux := http.ServeMux{}

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: handlers.LoggingHandler(os.Stdout, &serveMux),
	}

	s := Server{
		server:         &server,
		sessionManager: sessionManager,
	}

	serveMux.HandleFunc("/", s.route)

	return &s
}

func (s *Server) createSession(w http.ResponseWriter, r *http.Request, language string) {
	if r.Method != http.MethodGet {
		handleBadRequest(w, r, fmt.Errorf("unsupported method: %#+v", r.Method))
		return
	}

	session, err := s.sessionManager.CreateSession(language)
	if err != nil {
		handleBadRequest(w, r, err) // TODO may not always be a bad request; need typed errors
		return
	}

	handleCreateSessionResponse(w, r, session)
}

func (s *Server) pushToSession(w http.ResponseWriter, r *http.Request, rawSessionUUID string) {
	if r.Method != http.MethodPost {
		handleBadRequest(w, r, fmt.Errorf("unsupported method: %#+v", r.Method))
		return
	}

	pushToSessionRequest, err := handlePushToSessionRequest(w, r)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	sessionUUID, err := uuid.Parse(rawSessionUUID)
	if err != nil {
		handleBadRequest(w, r, fmt.Errorf("invalid session UUID: %#+v", rawSessionUUID))
		return
	}

	session, err := s.sessionManager.GetSession(sessionUUID)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	err = session.PushToSession(pushToSessionRequest.Data)
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	handlePushToSessionResponse(w, r, session)

	log.Printf(
		"!!! %v %v -> %#+v",
		r.Method,
		r.URL.Path,
		pushToSessionRequest.Data,
	)
}

func (s *Server) heartbeatForSession(w http.ResponseWriter, r *http.Request, rawSessionUUID string) {
	if r.Method != http.MethodGet {
		handleBadRequest(w, r, fmt.Errorf("unsupported method: %#+v", r.Method))
		return
	}

	sessionUUID, err := uuid.Parse(rawSessionUUID)
	if err != nil {
		handleBadRequest(w, r, fmt.Errorf("invalid session UUID: %#+v", rawSessionUUID))
		return
	}

	session, err := s.sessionManager.GetSession(sessionUUID)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	session.Heartbeat()

	handleHeartbeatForSessionResponse(w, r, session)

	log.Printf(
		"!!! %v %v -> heartbeat",
		r.Method,
		r.URL.Path,
	)
}

func (s *Server) proxySession(w http.ResponseWriter, r *http.Request, rawSessionUUID string) {
	sessionUUID, err := uuid.Parse(rawSessionUUID)
	if err != nil {
		handleBadRequest(w, r, fmt.Errorf("invalid session UUID: %#+v", rawSessionUUID))
		return
	}

	session, err := s.sessionManager.GetSession(sessionUUID)
	if err != nil {
		handleBadRequest(w, r, err)
		return
	}

	proxyUrl, err := url.Parse(session.GetProxyURL(r.URL))
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	log.Printf(
		"!!! %v %v -> %v",
		r.Method,
		r.URL.Path,
		proxyUrl.String(),
	)

	proxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.URL = proxyUrl

			request.Header.Set(
				"Host",
				fmt.Sprintf("%v:%v", session.Host(), session.Port()),
			)
		},
	}

	proxy.ServeHTTP(w, r)
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf(
		"<<< %v %v",
		r.Method,
		r.URL.Path,
	)

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 {
		handleBadRequest(w, r, fmt.Errorf("unknown URL path: %#+v", r.URL.Path))
		return
	}

	if len(parts) == 2 {
		if parts[0] == "create_session" {
			s.createSession(w, r, parts[1])
			return
		}

		if parts[0] == "push_to_session" {
			s.pushToSession(w, r, parts[1])
			return
		}

		if parts[0] == "heartbeat_for_session" {
			s.heartbeatForSession(w, r, parts[1])
			return
		}
	}

	if len(parts) >= 2 {
		if parts[0] == "proxy_session" {
			s.proxySession(w, r, parts[1])
			return
		}
	}

	handleBadRequest(w, r, fmt.Errorf("unknown URL path: %#+v", r.URL.Path))
}

func (s *Server) Open() error {
	var err error

	go func() {
		err = s.server.ListenAndServe()
	}()

	runtime.Gosched()

	time.Sleep(time.Millisecond)

	return err
}

func (s *Server) Close() {
	_ = s.server.Close()
}
