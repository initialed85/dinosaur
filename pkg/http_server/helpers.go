package http_server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/initialed85/dinosaur/pkg/sessions"
)

func getErrorResponseJSON(err error) ([]byte, error) {
	if err == nil {
		return []byte{}, fmt.Errorf("err cannot be nil")
	}

	response := ErrorResponse{Error: err.Error()}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return []byte{}, err
	}

	return responseJSON, nil
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, errToHandle error) {
	status := http.StatusBadRequest
	responseJSON, err := getErrorResponseJSON(errToHandle)
	if err != nil {
		status = http.StatusInternalServerError
		responseJSON = unknownInternalServerErrorResponseJSON
	}

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handleInternalServerError(w http.ResponseWriter, r *http.Request, errToHandle error) {

	status := http.StatusInternalServerError
	responseJSON, err := getErrorResponseJSON(errToHandle)
	if err != nil {
		responseJSON = unknownInternalServerErrorResponseJSON
	}

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handleGetSupportedLanguagesResponse(w http.ResponseWriter, r *http.Request, supportedLanguages []sessions.SupportedLanguage) {
	status := http.StatusOK
	responseJSON, err := json.Marshal(supportedLanguages)
	if err != nil {
		status = http.StatusInternalServerError
		responseJSON = unknownInternalServerErrorResponseJSON
	}

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handleCreateSessionResponse(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	response := CreateSessionResponse{
		UUID:        s.UUID(),
		Port:        s.Port(),
		InternalURL: s.InternalURL(),
		Code:        strings.TrimRight(strings.TrimLeft(s.Code(), "\r\n\t "), "\r\n\t ") + "\n",
	}

	status := http.StatusCreated
	responseJSON, err := json.Marshal(response)
	if err != nil {
		status = http.StatusInternalServerError
		responseJSON = unknownInternalServerErrorResponseJSON
	}

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handleGetSessionResponse(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	response := CreateSessionResponse{
		UUID:        s.UUID(),
		Port:        s.Port(),
		InternalURL: s.InternalURL(),
		Code:        strings.TrimRight(strings.TrimLeft(s.Code(), "\r\n\t "), "\r\n\t ") + "\n",
	}

	status := http.StatusOK
	responseJSON, err := json.Marshal(response)
	if err != nil {
		status = http.StatusInternalServerError
		responseJSON = unknownInternalServerErrorResponseJSON
	}

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handlePushToSessionRequest(w http.ResponseWriter, r *http.Request) (*PushToSessionRequest, error) {
	requestJSON, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	request := PushToSessionRequest{}

	err = json.Unmarshal(requestJSON, &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func handlePushToSessionResponse(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	status := http.StatusAccepted
	responseJSON := []byte(`{"success": true"}`)

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}

func handleHeartbeatForSessionResponse(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	status := http.StatusAccepted
	responseJSON := []byte(`{"success": true"}`)

	log.Printf(
		">>> %v %v %v %v",
		r.Method,
		r.URL.Path,
		status,
		string(responseJSON),
	)

	w.WriteHeader(status)
	_, _ = w.Write(responseJSON)
}
