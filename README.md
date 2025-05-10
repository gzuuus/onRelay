# onRely

A Go library providing atomic, high-performance ephemeral event storage and a set of flexible, composable policies for Nostr relays and related event-driven systems.

## Overview

`onRely` enables plug-and-play, modular storage and policy logic that can adapt to any relay's architecture. It consists of two main components:

1. **atomic**: Lock-free, fast ephemeral event storage
2. **policies**: Composable filter and event policy functions

## Installation

```bash
go get -u github.com/gzuuus/onRely
```

## Usage

### Atomic Storage

The `atomic` package provides a lock-free, fixed-size, thread-safe buffer for ephemeral event storage:

```go
import (
    "context"
    "github.com/gzuuus/onRely/atomic"
    "github.com/nbd-wtf/go-nostr"
)

// Create a new buffer with capacity for 1000 events
buffer := atomic.NewAtomicCircularBuffer2(1000)

// Save an event
err := buffer.SaveEvent(context.Background(), event)

// Query events
events, err := buffer.QueryEvents(context.Background(), filter)
```

### Policies

The `policies` package provides reusable, stateless functions for validating and modifying events and filters:

```go
import (
    "github.com/gzuuus/onRely/policies"
)

// Using filter policies in a relay
relay.RejectFilter = append(relay.RejectFilter,
    policies.NoEmptyFilters,
    policies.NoComplexFilters,
)

// Using event policies in a relay
relay.RejectEvent = append(relay.RejectEvent,
    policies.RejectEventsWithBase64Media,
    policies.PreventLargeTags(100),
    policies.RestrictToSpecifiedKinds(true, 1, 4, 5, 7),
)
```

## Design Philosophy

- No global aggregation: consumers decide how to combine, sequence, or gate policy evaluation
- Maximum flexibility and explicit modularity for advanced and evolving relay architectures
- Standard reason strings allow app-layer protocol signaling

## License

[MIT License](LICENSE)
