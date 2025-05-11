# onRelay Policies

This package provides policy functions for filtering and validating Nostr events and filters.

## Cross-Framework Compatibility

The policies in this package are designed to be framework-agnostic, allowing them to be used with various Nostr relay implementations. This is achieved through a simple dual-layer pattern:

1. **Core Logic**: All policy functions return a canonical `error` type.
2. **Adapters**: Simple adapter functions convert error responses into framework-specific signatures.

## Using Policies

### Filter Policies

Filter policies validate incoming filter requests and reject those that don't meet certain criteria:

```go
// Using a filter policy directly (error-based)
err := policies.NoComplexFilters(ctx, filter)
if err != nil {
    // Handle rejection
}

// Using with a framework that expects (bool, string)
rejectFilter := policies.ToBoolStringFilter(policies.NoComplexFilters)
reject, reason := rejectFilter(ctx, filter)
if reject {
    // Handle rejection with reason
}
```

### Event Policies

Event policies validate incoming events and reject those that don't meet certain criteria:

```go
// Using an event policy directly (error-based)
err := policies.RejectEventsWithBase64Media(ctx, event)
if err != nil {
    // Handle rejection
}

// Using with a framework that expects (bool, string)
rejectEvent := policies.ToBoolStringEventFilter(policies.RejectEventsWithBase64Media)
reject, reason := rejectEvent(ctx, event)
if reject {
    // Handle rejection with reason
}
```

## Framework Integration

### For Error-Based Frameworks (e.g., rely)

```go
import (
    "github.com/your-org/onRelay/policies"
)

// Direct integration
rely.RejectFilters = append(rely.RejectFilters, policies.NoComplexFilters)
rely.RejectEvents = append(rely.RejectEvents, policies.RejectEventsWithBase64Media)

// Using combined policies
combinedFilter := policies.CombineFilters(policies.NoComplexFilters, policies.NoEmptyFilters)
rely.RejectFilters = append(rely.RejectFilters, combinedFilter)
```

### For (bool, string)-Based Frameworks (e.g., khatru)

```go
import (
    "github.com/your-org/onRelay/policies"
)

// Using adapters for single policies
khatru.RejectFilter = append(
    khatru.RejectFilter,
    policies.ToBoolStringFilter(policies.NoComplexFilters),
)
khatru.RejectEvent = append(
    khatru.RejectEvent,
    policies.ToBoolStringEventFilter(policies.RejectEventsWithBase64Media),
)

// Using adapters for combined policies
combinedFilter := policies.CombineFilters(policies.NoComplexFilters, policies.NoEmptyFilters)
khatru.RejectFilter = append(
    khatru.RejectFilter,
    policies.ToBoolStringFilter(combinedFilter),
)
```

## Adapting to Different Frameworks

The core design allows you to adapt policies to any framework's expected return type. Here are examples for common scenarios:

### Custom Return Types

For frameworks with custom return types, you can create your own adapters:

```go
// Example: Framework expects a custom response type
type CustomResponse struct {
    Allowed bool
    Reason  string
    Code    int
}

// Create an adapter for the custom response type
func ToCustomResponseFilter(f policies.FilterPolicy) func(ctx context.Context, filter nostr.Filter) CustomResponse {
    return func(ctx context.Context, filter nostr.Filter) CustomResponse {
        err := f(ctx, filter)
        if err != nil {
            return CustomResponse{Allowed: false, Reason: err.Error(), Code: 400}
        }
        return CustomResponse{Allowed: true, Reason: "", Code: 200}
    }
}

// Usage
customFilter := ToCustomResponseFilter(policies.NoComplexFilters)
response := customFilter(ctx, filter)
```

### HTTP Handler Integration

Adapt policies to work with HTTP handlers:

```go
// Create an HTTP handler that uses a filter policy
func FilterPolicyHandler(policy policies.FilterPolicy) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract filter from request
        var filter nostr.Filter
        if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
            http.Error(w, "Invalid filter", http.StatusBadRequest)
            return
        }
        
        // Apply policy
        if err := policy(r.Context(), filter); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Filter accepted"))
    }
}

// Usage
http.HandleFunc("/filter", FilterPolicyHandler(policies.NoComplexFilters))
```

### Creating Custom Policies

To create a new policy, simply define a function that returns an error:

```go
func MyCustomFilterPolicy(ctx context.Context, filter nostr.Filter) error {
    // Your validation logic here
    if someCondition {
        return errors.New("rejection reason")
    }
    return nil
}
```
