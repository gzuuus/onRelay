package atomic

import (
	"context"

	"github.com/gzuuus/onRelay"

	"github.com/nbd-wtf/go-nostr"
)

// ReplaceEvent is not implemented for AtomicCircularBuffer as it's designed for ephemeral storage.
func (cb *AtomicCircularBuffer) ReplaceEvent(ctx context.Context, evt *nostr.Event) error {
	return onRelay.ErrNotImplemented
}

// DeleteEvent is not implemented for AtomicCircularBuffer as it's designed for ephemeral storage.
func (cb *AtomicCircularBuffer) DeleteEvent(ctx context.Context, evt *nostr.Event) error {
	return onRelay.ErrNotImplemented
}
