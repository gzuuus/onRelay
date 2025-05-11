package policies

import (
	"context"
	"errors"

	"github.com/nbd-wtf/go-nostr"
)

// FilterPolicy defines the core function type for filter validation policies.
// It returns an error if the filter should be rejected, or nil if it should be accepted.
type FilterPolicy func(ctx context.Context, filter nostr.Filter) error

// NoComplexFilters disallows filters with more than 2 tags.
func NoComplexFilters(ctx context.Context, filter nostr.Filter) error {
	items := len(filter.Tags) + len(filter.Kinds)

	if items > 4 && len(filter.Tags) > 2 {
		return errors.New("too many things to filter for")
	}

	return nil
}

// NoEmptyFilters disallows filters that don't have at least a tag, a kind, an author or an id.
func NoEmptyFilters(ctx context.Context, filter nostr.Filter) error {
	c := len(filter.Kinds) + len(filter.IDs) + len(filter.Authors)
	for _, tagItems := range filter.Tags {
		c += len(tagItems)
	}
	if c == 0 {
		return errors.New("can't handle empty filters")
	}
	return nil
}