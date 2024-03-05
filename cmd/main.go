package main

import (
	"log"

	"github.com/initialed85/dinosaur/internal/cmd"
	"github.com/initialed85/dinosaur/pkg/http_server"
	"github.com/initialed85/dinosaur/pkg/sessions"
)

func main() {
	sessionManager := sessions.NewManager()

	server := http_server.New(
		8080,
		sessionManager,
	)

	err := sessionManager.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sessionManager.Close()

	err = server.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	cmd.WaitForSIGTERM()
}
