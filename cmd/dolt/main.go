package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello %#v", *r)
	w.Write([]byte("OK"))
}

func main() {
	// TODO: Read configuration from environment

	sigint := make(chan os.Signal, 1)

	// TODO: Define list of ignored syscall Signals
	// signal.Ignore(syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)

	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGQUIT)

	healthSrv := http.Server{
		Addr:    ":8080",
		Handler: &Handler{},
	}

	go func() {
		<-sigint
		// TODO: Postpone server shutdown regarding Env variable
		if err := healthSrv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown with error: %v", err)
		}
	}()

	// TODO: Postpone server start regarding Env variable

	log.Printf("Server has been started and listen on: %s", healthSrv.Addr)
	if err := healthSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	// TODO: Suden interrupt with non-0 exit code
	log.Println("Server has been stopped")
}
