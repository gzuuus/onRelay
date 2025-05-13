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

onRelay policies are composable and can be easily adapted to various frameworks using the provided adapters. See the [examples directory](./examples) for up-to-date usage.

#### Example: Combining and Adapting Policies for rely

See [`examples/combine.go`](./examples/combine.go) for how to combine multiple event policies and use them with `ToClientAdapter` for frameworks like rely:

```go
import (
    "github.com/gzuuus/onRelay/policies"
    "github.com/pippellia-btc/rely"
)

combinedEventPolicy := policies.CombineEvents(
    policies.RestrictToSpecifiedKinds(1, 30023),
    policies.RejectEventsWithBase64Media,
)
rely.RejectEvent = append(rely.RejectEvent, policies.ToClientAdapter[rely.Client](combinedEventPolicy))
```

#### Example: Adapting QueryEvents for khatru

See [`examples/khatru.go`](./examples/khatru.go) for how to use the `QueryEventsToChan` adapter to make onRelay's storage compatible with khatru's channel-based API:

```go
import (
    "github.com/fiatjaf/khatru"
    "github.com/gzuuus/onRelay/atomic"
    "github.com/nbd-wtf/go-nostr"
)

// The QueryEventsToChan adapter is provided in the atomic package
// It converts slice-based event queries to channel-based ones for frameworks like khatru

// Usage with khatru:
var buffer = atomic.NewAtomicCircularBuffer(1000)
relay := khatru.NewRelay()

// Store events in the buffer
relay.StoreEvent = append(relay.StoreEvent, buffer.SaveEvent)

// Use the adapter to convert QueryEvents to the channel-based signature khatru expects
relay.QueryEvents = append(relay.QueryEvents, atomic.QueryEventsToChan(buffer.QueryEvents))

// Full example available in examples/khatru.go
```

For more, browse the [examples](./examples) folder.

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
