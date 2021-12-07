package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	exitCode   int
	healthUri  url.URL
	ignoreSigs []os.Signal
	initTime   time.Duration
	lifeTime   time.Duration
	stopTime   time.Duration
}

func newConfig() *config {
	return &config{
		exitCode: 0,
		healthUri: url.URL{
			Host: ":8080",
		},
		ignoreSigs: []os.Signal{},
		initTime:   0,
		lifeTime:   0,
		stopTime:   0,
	}
}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	cfg := newConfig()

	sigint := make(chan os.Signal, 1)

	if len(cfg.ignoreSigs) > 0 {
		signal.Ignore(cfg.ignoreSigs...)
	}

	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGQUIT)

	healthSrv := http.Server{
		Addr:    cfg.healthUri.Host,
		Handler: &Handler{},
	}

	go func() {
		<-sigint
		if cfg.stopTime > 0 {
			<-time.After(cfg.stopTime)
		}
		if err := healthSrv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown with error: %v", err)
		}
	}()

	// TODO: Postpone server start regarding Env variable
	if cfg.initTime > 0 {
		<-time.After(cfg.initTime)
	}
	log.Printf("Server has been started and listen on: %s", healthSrv.Addr)
	if err := healthSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	// TODO: Suden interrupt with non-0 exit code
	log.Println("Server has been stopped")
}
