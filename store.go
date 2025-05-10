package onRely

import (
	"context"
	"errors"

	"github.com/nbd-wtf/go-nostr"
)

// Common errors that may be returned by Store implementations
var (
	// ErrNotImplemented is returned for optional operations that are not implemented.
	ErrNotImplemented = errors.New("operation not implemented")

	// ErrNilEvent is returned when a nil event is passed to SaveEvent or ReplaceEvent.
	ErrNilEvent = errors.New("event cannot be nil")
)

// Store is a composable event storage interface.
// It defines the core operations for storing and retrieving Nostr events.
type Store interface {
	// SaveEvent adds a new event to the store.
	SaveEvent(ctx context.Context, evt *nostr.Event) error

	// QueryEvents returns a slice of events matching the filter.
	QueryEvents(ctx context.Context, filter nostr.Filter) ([]*nostr.Event, error)

	// Optional operations may return ErrNotImplemented if not supported
	
	// ReplaceEvent replaces an existing event with the same ID.
	ReplaceEvent(ctx context.Context, evt *nostr.Event) error

	// DeleteEvent removes an event from the store.
	DeleteEvent(ctx context.Context, evt *nostr.Event) error

	// CountEvents returns the number of events matching the filter.
	CountEvents(ctx context.Context, filter nostr.Filter) (int, error)
}
