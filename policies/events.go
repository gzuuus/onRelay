package policies

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

// EventPolicy defines a function type for event validation policies.
// It returns a boolean indicating if the event should be rejected and a reason message.
type EventPolicy func(context.Context, *nostr.Event) (bool, string)

// RestrictToSpecifiedKinds returns a function that can be used as a RejectFilter that will reject
// any events with kinds different than the specified ones.
func RestrictToSpecifiedKinds(kinds ...uint16) EventPolicy {
	// sort the kinds in increasing order
	slices.Sort(kinds)

	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if _, allowed := slices.BinarySearch(kinds, uint16(event.Kind)); allowed {
			return false, ""
		}

		return true, fmt.Sprintf("received event kind %d not allowed", event.Kind)
	}
}

// RestrictToSpecifiedKindsRanges returns a policy function that restricts events based on their kind ranges.
// It allows you to specify whether to allow ephemeral events, regular events, and replaceable events.
func RestrictToSpecifiedKindsRanges(allowEphemeral bool, allowRegular bool, allowReplaceable bool) EventPolicy {
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		// Check if it's an ephemeral event (kind 20000-29999)
		if nostr.IsEphemeralKind(event.Kind) {
			return !allowEphemeral, "ephemeral events are not allowed"
		}

		// Check if it's a replaceable event (kind 0, 1, 2, 3, 10000-19999, 30000-39999)
		if nostr.IsReplaceableKind(event.Kind) {
			return !allowReplaceable, "replaceable events are not allowed"
		}

		if nostr.IsRegularKind(event.Kind) {
			return !allowRegular, "regular events are not allowed"
		}

		return false, ""
	}
}

// PreventTimestampsInThePast rejects events with timestamps older than the specified threshold.
func PreventTimestampsInThePast(threshold time.Duration) EventPolicy {
	thresholdSeconds := nostr.Timestamp(threshold.Seconds())
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if nostr.Now()-event.CreatedAt > thresholdSeconds {
			return true, "event too old"
		}
		return false, ""
	}
}

// PreventTimestampsInTheFuture rejects events with timestamps too far in the future.
func PreventTimestampsInTheFuture(threshold time.Duration) EventPolicy {
	thresholdSeconds := nostr.Timestamp(threshold.Seconds())
	return func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
		if event.CreatedAt-nostr.Now() > thresholdSeconds {
			return true, "event too much in the future"
		}
		return false, ""
	}
}

// RejectEventsWithBase64Media rejects events containing base64-encoded media.
func RejectEventsWithBase64Media(ctx context.Context, evt *nostr.Event) (bool, string) {
	return strings.Contains(evt.Content, "data:image/") || strings.Contains(evt.Content, "data:video/"), "event with base64 media"
}
