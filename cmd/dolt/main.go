package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	cfg := &config{
		exitCode: 0,
		healthUri: url.URL{
			Scheme: "http",
			Host:   ":8080",
			Path:   "/health",
		},
		ignoreSigs: []os.Signal{},
		initTime:   0,
		lifeTime:   0,
		stopTime:   0,
	}

	if val := os.Getenv("EXITCODE"); val != "" {
		if val, err := strconv.Atoi(val); err == nil {
			cfg.exitCode = val
		}
	}

	if val := os.Getenv("HEALTHURI"); val != "" {
		if val, err := url.Parse(val); err == nil {
			cfg.healthUri = *val
		}
	}

	if vals := os.Getenv("IGNORESIGS"); vals != "" {
		for _, v := range strings.Split(vals, ",") {
			if val, err := strconv.Atoi(v); err == nil {
				cfg.ignoreSigs = append(cfg.ignoreSigs, syscall.Signal(val))
			}
		}
	}

	if val := os.Getenv("INITTIME"); val != "" {
		if val, err := time.ParseDuration(val); err == nil {
			cfg.initTime = val
		}
	}

	if val := os.Getenv("LIFETIME"); val != "" {
		if val, err := time.ParseDuration(val); err == nil {
			cfg.initTime = val
		}
	}

	if val := os.Getenv("STOPTIME"); val != "" {
		if val, err := time.ParseDuration(val); err == nil {
			cfg.initTime = val
		}
	}

	return cfg
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

	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	if cfg.lifeTime > 0 {
		go func() {
			<-time.After(cfg.lifeTime)
			close(sigint)
		}()
	}

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
		log.Println("Server has been stopped")
	}()

	if cfg.initTime > 0 {
		<-time.After(cfg.initTime)
	}
	log.Printf("Server has been started and listen on: %s", healthSrv.Addr)
	if err := healthSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	os.Exit(cfg.exitCode)
}
