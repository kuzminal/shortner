package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"short/internal/httpio"
	"short/short"
	"time"
)

func main() {
	const (
		addr    = "localhost:8080"
		timeout = 10 * time.Second
	)
	fmt.Fprintln(os.Stderr, "starting the server on", addr)

	shortener := short.NewServer()

	logger := log.New(os.Stderr, "shortener: ", log.LstdFlags|log.Lmsgprefix)
	logger.Println("starting the server on", addr)

	server := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(shortener, timeout, "timeout"),
		ReadTimeout: timeout,
	}

	if os.Getenv("LINKIT_DEBUG") == "1" {
		server.ErrorLog = logger
		server.Handler = httpio.LoggingMiddleware(server.Handler)
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Println("server closed unexpectedly:", err)
	}
}
