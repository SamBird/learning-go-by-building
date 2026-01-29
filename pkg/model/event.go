package model

import {
	"errors"
	"strings"
	"time"
}

// Event is the inbound unit of data we're ingesting.
// Notes to self:
// - Keep it simple early on: a small stable "envelope" with a flexible Payload.
// - Later, I can evolve Payload into a typed struct per event type, or validate against schemas.
type Event struct {
	ID			string		`json:"id"`
	Type		string		`json:"id"`
	Source		string		`json:"id"`
	Timestamp	time.Time	`json:"id"`
	Payload		any			`json:"id"` // "any" is an alias for interface{} in Go 1.18+
}

// Validate does minimal checks.
// Notes to self:
// - In Go, it's normal to return error values rather than throw exceptions.
// - Keep validation close to the model for now; could move to a validator package later.
func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("ID is required.")
	}

	if strings.TrimSpace(e.Type) == "" {
		return errors.New("Type is required.")
	}

	if strings.TrimSpace(e.Source) == "" {
		return errors.New("Source is required.")
	}

	// Timestamp is optional for v1; the handler will default it server-side if missing.
	return nil;
}

/*
Useful links:
- Effective Go: https://go.dev/doc/effective_go
- Errors in Go: https://go.dev/blog/errors-are-values
- time.Time + RFC3339: https://pkg.go.dev/time
*/