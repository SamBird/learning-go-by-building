package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/SamBird/learning-go-by-building/pkg/model"
)

// EventHandler owns HTTP handlers + dependencies.
//
// Notes to self:
// - This is "poor man's DI": pass dependencies in via a constructor.
// - It keeps things testable without bringing in a framework.
// https://gobyexample.com/pointers
type EventHandler struct {
	Logger *log.Logger
}

func NewEventHandler(logger *log.Logger) *EventHandler {
	return &EventHandler{Logger: logger}
}

// Register wires endpoints into a ServeMux.
//
// Notes to self:
// - net/http + ServeMux is perfectly fine for a lot of services.
// - The "METHOD /path" patterns are supported in newer Go versions (Go 1.22+).
func (h *EventHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /events", h.handlePostEvent)
	mux.HandleFunc("GET /health", h.handleHealth)
}

// handleHealth is a simple liveness endpoint.
// Notes to self:
// - In real deployments I'd normally separate /health (liveness) and /ready (readiness).
func (h *EventHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

