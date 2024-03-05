package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSIGTERM() {
	c := make(chan os.Signal, 16)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
