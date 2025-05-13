package atomic

import (
	"context"
)

func SliceToChanEvents(fn func(ctx context.Context, filter nostr.Filter) ([]*nostr.Event, error)) func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
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
