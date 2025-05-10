package atomic

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

// CountEvents returns the number of events matching the filter.
func (cb *AtomicCircularBuffer) CountEvents(ctx context.Context, filter nostr.Filter) (int, error) {
	events, err := cb.QueryEvents(ctx, filter)
	if err != nil {
		return 0, err
	}
	return len(events), nil
}
