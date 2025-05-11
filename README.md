# onRelay

A framework-agnostic Go library providing plug and play event storage and adaptable policy functions for Nostr relays.

## Features

- **Cross-Framework Compatibility**: Core policies return canonical `error` types with adapters for any signature format
- **Atomic Storage**: Lock-free circular buffer optimized for high-throughput ephemeral event storage
- **Composable Policies**: Modular filter and event validation functions that can be combined as needed

## Installation

```bash
go get -u github.com/gzuuus/onRelay
```

## Usage

### Atomic Storage

```go
import (
    "context"
    "github.com/gzuuus/onRelay/atomic"
    "github.com/nbd-wtf/go-nostr"
)

func main() {
    buffer := atomic.NewCircularBuffer(10_000)
    buffer.SaveEvent(context.Background(), &nostr.Event{...})
    events := buffer.QueryEvents(context.Background(), nostr.Filter{...})
}
```

### Framework-Agnostic Policies

```go
import (
    "github.com/gzuuus/onRelay/policies"
)

// For error-based frameworks (e.g., rely)
rely.RejectFilters = append(rely.RejectFilters, policies.NoComplexFilters)

// For (bool, string)-based frameworks (e.g., khatru)
khatru.RejectFilter = append(
    khatru.RejectFilter,
    policies.ToBoolStringFilter(policies.NoComplexFilters),
)

// Combining multiple policies
combinedFilter := policies.CombineFilters(policies.NoComplexFilters, policies.NoEmptyFilters)

// Combined policies with adapter
policies.ToBoolStringFilter(combinedFilter)
```

### Custom Framework Adapters

```go
// Adapt to any custom framework response type
func ToCustomResponse(f policies.FilterPolicy) func(ctx, filter) CustomType {
    return func(ctx context.Context, filter nostr.Filter) CustomType {
        err := f(ctx, filter)
        if err != nil {
            return CustomType{Rejected: true, Reason: err.Error()}
        }
        return CustomType{Rejected: false}
    }
}
```

## Design Philosophy

- **Framework Agnostic**: Core logic is decoupled from framework-specific signatures
- **Minimalist API**: Focused on essential functionality without unnecessary abstractions
- **Composable Units**: Small, single-purpose functions that can be combined as needed
- **Explicit Control**: Decide how to combine, sequence, or adapt policy evaluation

## License

[MIT License](LICENSE)
