package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/<your-github-username>/learning-go-by-building/internal/handler"
)

func main() {
	// Notes to self:
	// - log.New lets me control output + prefix/flags.
	// - LUTC prints timestamps in UTC which is usually better for systems logs.
	// Ref: https://pkg.go.dev/log
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	// ServeMux is the built-in router.
	// Notes to self:
	// - I’m intentionally using the standard library to learn Go properly.
	// - Frameworks are great later; stdlib teaches the fundamentals.
	mux := http.NewServeMux()

	h := handler.NewEventHandler(logger)
	h.Register(mux)

	// Configure HTTP server with a basic timeout.
	// Notes to self:
	// - ReadHeaderTimeout protects against slowloris-style attacks.
	// - Later: add ReadTimeout/WriteTimeout/IdleTimeout once I understand their tradeoffs.
	// Ref: https://pkg.go.dev/net/http#Server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("starting server on %s", srv.Addr)

	// ListenAndServe blocks.
	// Notes to self:
	// - In the next iteration, add graceful shutdown (context + signal handling).
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server failed: %v", err)
	}
}

/*
Useful links (notes to self):
- HTTP server basics: https://pkg.go.dev/net/http
- Go Proverbs (idioms): https://go-proverbs.github.io/
- Effective Go (general style): https://go.dev/doc/effective_go
*/
