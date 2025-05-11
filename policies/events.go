package policies

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

// EventPolicy defines the core function type for event validation policies.
// It returns an error if the event should be rejected, or nil if it should be accepted.
type EventPolicy func(context.Context, *nostr.Event) error

// RestrictToSpecifiedKinds returns a function that will reject any events with kinds
// different than the specified ones.
func RestrictToSpecifiedKinds(kinds ...uint16) EventPolicy {
	// sort the kinds in increasing order
	slices.Sort(kinds)

	return func(ctx context.Context, event *nostr.Event) error {
		if _, allowed := slices.BinarySearch(kinds, uint16(event.Kind)); allowed {
			return nil
		}

		return fmt.Errorf("received event kind %d not allowed", event.Kind)
	}
}

// RestrictToSpecifiedKindsRanges returns a policy function that restricts events based on their kind ranges.
// It allows you to specify whether to allow ephemeral events, regular events, and replaceable events.
func RestrictToSpecifiedKindsRanges(allowEphemeral bool, allowRegular bool, allowReplaceable bool) EventPolicy {
	return func(ctx context.Context, event *nostr.Event) error {
		if nostr.IsEphemeralKind(event.Kind) && !allowEphemeral {
			return errors.New("ephemeral events are not allowed")
		}

		if nostr.IsRegularKind(event.Kind) && !allowRegular {
			return errors.New("regular events are not allowed")
		}

		if nostr.IsReplaceableKind(event.Kind) && !allowReplaceable {
			return errors.New("replaceable events are not allowed")
		}

		return nil
	}
}

// PreventTimestampsInThePast rejects events with timestamps older than the specified threshold.
func PreventTimestampsInThePast(threshold time.Duration) EventPolicy {
	thresholdSeconds := nostr.Timestamp(threshold.Seconds())
	return func(ctx context.Context, event *nostr.Event) error {
		if nostr.Now()-event.CreatedAt > thresholdSeconds {
			return errors.New("event too old")
		}
		return nil
	}
}

// PreventTimestampsInTheFuture rejects events with timestamps too far in the future.
func PreventTimestampsInTheFuture(threshold time.Duration) EventPolicy {
	thresholdSeconds := nostr.Timestamp(threshold.Seconds())
	return func(ctx context.Context, event *nostr.Event) error {
		if event.CreatedAt-nostr.Now() > thresholdSeconds {
			return errors.New("event too much in the future")
		}
		return nil
	}
}

// RejectEventsWithBase64Media rejects events containing base64-encoded media.
func RejectEventsWithBase64Media(ctx context.Context, evt *nostr.Event) error {
	if strings.Contains(evt.Content, "data:image/") || strings.Contains(evt.Content, "data:video/") {
		return errors.New("event with base64 media")
	}
	return nil
}
