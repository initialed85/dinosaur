package main

import (
	"github.com/initialed85/dinosaur/pkg/http_server"
	"github.com/initialed85/dinosaur/pkg/sessions"
	"log"
)

func main() {
	m := sessions.NewManager()
	defer m.Close()

	s := http_server.New(
		8080,
		m,
	)

	err := s.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()
}
