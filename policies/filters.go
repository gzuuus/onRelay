package policies

import (
	"context"

	"github.com/nbd-wtf/go-nostr"
)

// FilterPolicy defines a function type for filter validation policies.
// It returns a boolean indicating if the filter should be rejected and a reason message.
type FilterPolicy func(ctx context.Context, filter nostr.Filter) (reject bool, msg string)

// NoComplexFilters disallows filters with more than 2 tags.
func NoComplexFilters(ctx context.Context, filter nostr.Filter) (reject bool, msg string) {
	items := len(filter.Tags) + len(filter.Kinds)

	if items > 4 && len(filter.Tags) > 2 {
		return true, "too many things to filter for"
	}

	return false, ""
}

// NoEmptyFilters disallows filters that don't have at least a tag, a kind, an author or an id.
func NoEmptyFilters(ctx context.Context, filter nostr.Filter) (reject bool, msg string) {
	c := len(filter.Kinds) + len(filter.IDs) + len(filter.Authors)
	for _, tagItems := range filter.Tags {
		c += len(tagItems)
	}
	if c == 0 {
		return true, "can't handle empty filters"
	}
	return false, ""
}
