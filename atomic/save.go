package atomic

import (
	"context"

	"github.com/gzuuus/onRelay"

	"github.com/nbd-wtf/go-nostr"
)

// SaveEvent adds a new event to the circular buffer.
// If the buffer is full, it automatically overwrites the oldest event.
func (cb *AtomicCircularBuffer) SaveEvent(ctx context.Context, evt *nostr.Event) error {
	if evt == nil {
		return onRelay.ErrNilEvent
	}

	head := cb.head.Load()
	cb.buffer[head].Store(evt)
	cb.head.Store((head + 1) % cb.size)

	count := cb.count.Add(1)
	if count > cb.size {
		cb.count.Store(cb.size)
	}

	return nil
}
