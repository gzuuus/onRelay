package atomic

import (
	"context"
	
	"github.com/nbd-wtf/go-nostr"
)

// QueryEventsToChan converts a function that returns a slice of events to one that returns a channel of events
// This is useful for adapting to frameworks like khatru that expect channel-based event delivery
func QueryEventsToChan(fn func(ctx context.Context, filter nostr.Filter) ([]*nostr.Event, error)) func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
	return func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
		events, err := fn(ctx, filter)
		ch := make(chan *nostr.Event, len(events))
		if err != nil {
			close(ch)
			return ch, err
		}
		go func() {
			defer close(ch)
			for _, event := range events {
				ch <- event
			}
		}()
		return ch, nil
	}
}
