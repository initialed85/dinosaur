package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSigInt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	<-c
}
