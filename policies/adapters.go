package policies

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

// This file contains simple adapter functions to help with cross-framework compatibility.
// These adapters convert between different signature formats to enable seamless integration
// with various Nostr relay frameworks.

// CombineFilters combines multiple filter policies into a single policy.
// It returns an error if any of the policies returns an error.
func CombineFilters(policies ...FilterPolicy) FilterPolicy {
	return func(ctx context.Context, filter nostr.Filter) error {
		for _, policy := range policies {
			if err := policy(ctx, filter); err != nil {
				return err
			}
		}
		return nil
	}
}

// CombineEvents combines multiple event policies into a single policy.
// It returns an error if any of the policies returns an error.
func CombineEvents(policies ...EventPolicy) EventPolicy {
	return func(ctx context.Context, event *nostr.Event) error {
		for _, policy := range policies {
			if err := policy(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}
}

// ToBoolStringFilter converts a core error-based filter policy to a (bool, string) signature.
// This is useful for frameworks that expect the (bool, string) return format.
func ToBoolStringFilter(f FilterPolicy) func(ctx context.Context, filter nostr.Filter) (bool, string) {
	return func(ctx context.Context, filter nostr.Filter) (bool, string) {
		err := f(ctx, filter)
		if err != nil {
			return true, err.Error()
		}
		return false, ""
	}
}

// ToBoolStringEventFilter converts a core error-based event policy to a (bool, string) signature.
// This is useful for frameworks that expect the (bool, string) return format.
func ToBoolStringEventFilter(f EventPolicy) func(context.Context, *nostr.Event) (bool, string) {
	return func(ctx context.Context, event *nostr.Event) (bool, string) {
		err := f(ctx, event)
		if err != nil {
			return true, err.Error()
		}
		return false, ""
	}
}
