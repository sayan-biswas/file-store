package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sayan-biswas/file-store/pkg/server/server"
)

func main() {

	// start server in a separate Go routine
	go server.Start()

	// stop server on TERM singnal from OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
