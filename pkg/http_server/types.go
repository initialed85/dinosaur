package http_server

import (
	"github.com/google/uuid"
)

var (
	unknownInternalServerErrorResponseJSON = []byte(`{"error": "Unknown internal server error"}`)
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateSessionResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Port        int       `json:"port"`
	InternalURL string    `json:"internal_url"`
	Code        string    `json:"code"`
}

type PushToSessionRequest struct {
	Data string `json:"data"`
}
