package main

import (
	"github.com/initialed85/dinosaur/internal/cmd"
	"github.com/initialed85/dinosaur/pkg/http_server"
	"github.com/initialed85/dinosaur/pkg/sessions"
	"log"
)

func main() {
	s := http_server.New(
		8080,
		sessions.NewManager(),
	)

	go func() {
		err := s.Open()
		if err != nil {
			log.Fatal(err)
		}
	}()

	defer s.Close()

	cmd.WaitForSigInt()
}
