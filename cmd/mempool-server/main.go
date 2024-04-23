package main

import (
	"log/slog"
	"net/http"
	"time"

	"mempool-server/pkg/routes"
)

const localAddr = ":48080"

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	routes.Register(mux)

	server := http.Server{
		Addr:              localAddr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
	}
	slog.Info("Starting server on", slog.String("addr", localAddr))
	return server.ListenAndServe()
}
