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
//   If this ever fails in older versions, switch to mux.HandleFunc("/events", ...) + method checks.
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

// handlePostEvent accepts an event payload, validates it, and (for now) logs it.
func (h *EventHandler) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	// Notes to self:
	// - Always protect against huge payloads. MaxBytesReader prevents memory bloat / DoS-ish behaviour.
	// - 1<<20 = 1,048,576 bytes = 1MB (fine for v1).
	// Ref: https://pkg.go.dev/net/http#MaxBytesReader
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var evt model.Event

	// json.Decoder is streaming-friendly vs reading the whole body then Unmarshal.
	// Ref: https://pkg.go.dev/encoding/json#Decoder
	dec := json.NewDecoder(r.Body)

	// DisallowUnknownFields is useful to catch typos early (client sends "soucre" etc).
	// Caveat: it can be strict when payloads evolve; for v1 it's good discipline.
	// Ref: https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields
	dec.DisallowUnknownFields()

	// Decode into our struct (JSON -> Go struct).
	if err := dec.Decode(&evt); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON payload", err)
		return
	}

	// Notes to self:
	// - This is a small safety check to ensure the request isn't "valid JSON + extra junk".
	// - Helps avoid weird edge cases where multiple JSON values are sent.
	if dec.More() {
		httpError(w, http.StatusBadRequest, "unexpected extra JSON content", errors.New("multiple JSON values"))
		return
	}

	// Default timestamp server-side if it's missing.
	// Notes to self:
	// - UTC is a good default for server logs + events.
	// - RFC3339 is the common wire-format.
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now().UTC()
	}

	// Validate required fields.
	if err := evt.Validate(); err != nil {
		httpError(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Structured logging is the next step, but for now a consistent message format is fine.
	// Later: swap to slog (Go's structured logger) or zap/zerolog.
	h.Logger.Printf(
		"event accepted: id=%s type=%s source=%s ts=%s",
		evt.ID,
		evt.Type,
		evt.Source,
		evt.Timestamp.Format(time.RFC3339),
	)

	// Reply to the client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "accepted",
		"id":     evt.ID,
	})
}

// httpError returns a consistent JSON error shape.
// Notes to self:
// - Keep error responses consistent from day one; it helps clients + debugging.
// - "details" can leak internal info in real services; later we might hide it in prod mode.
func httpError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":   message,
		"details": err.Error(),
	})
}

/*
Useful links (notes to self):
- net/http package docs: https://pkg.go.dev/net/http
- ServeMux patterns (Go 1.22): https://go.dev/blog/routing-enhancements
- JSON decoding tips: https://pkg.go.dev/encoding/json
*/
