package atomic

import (
	"onRelay"
	"sync/atomic"

	"github.com/nbd-wtf/go-nostr"
)

var _ onRelay.Store = (*AtomicCircularBuffer)(nil)

// AtomicCircularBuffer is an optimized, lock-free, fixed-size circular buffer for storing Nostr events.
type AtomicCircularBuffer struct {
	buffer []*atomic.Pointer[nostr.Event]
	head   atomic.Uint64 // position to write next event
	size   uint64        // fixed size of the buffer
	count  atomic.Uint64 // number of events in buffer
}

// NewAtomicCircularBuffer creates a new AtomicCircularBuffer with the specified capacity.
func NewAtomicCircularBuffer(capacity int) *AtomicCircularBuffer {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	buffer := make([]*atomic.Pointer[nostr.Event], capacity)
	for i := range buffer {
		buffer[i] = &atomic.Pointer[nostr.Event]{}
	}

	return &AtomicCircularBuffer{
		buffer: buffer,
		size:   uint64(capacity),
	}
}
