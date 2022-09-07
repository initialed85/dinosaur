package http_server

import (
	"encoding/json"
	"fmt"
	"github.com/initialed85/dinosaur/pkg/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	status := http.StatusInternalServerError
	responseJSON := unknownInternalServerErrorResponseJSON
	var err error

	responseJSON, err = getErrorResponseJSON(errToHandle)
	if err == nil {
		status = http.StatusBadRequest
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
	responseJSON := unknownInternalServerErrorResponseJSON

	responseJSON, _ = getErrorResponseJSON(errToHandle)

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
	status := http.StatusInternalServerError
	responseJSON := unknownInternalServerErrorResponseJSON
	var err error

	responseJSON, err = json.Marshal(supportedLanguages)
	if err == nil {
		status = http.StatusOK
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
	status := http.StatusInternalServerError
	responseJSON := unknownInternalServerErrorResponseJSON
	var err error

	response := CreateSessionResponse{
		UUID:        s.UUID(),
		Port:        s.Port(),
		InternalURL: s.InternalURL(),
		Code:        strings.TrimRight(strings.TrimLeft(s.Code(), "\r\n\t "), "\r\n\t ") + "\n",
	}

	responseJSON, err = json.Marshal(response)
	if err == nil {
		status = http.StatusCreated
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
	requestJSON, err := ioutil.ReadAll(r.Body)
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
	responseJSON := unknownInternalServerErrorResponseJSON

	// TODO all of this
	responseJSON = []byte(`{"success": true"}`)

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
	responseJSON := unknownInternalServerErrorResponseJSON

	// TODO all of this
	responseJSON = []byte(`{"success": true"}`)

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
