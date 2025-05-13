package main

import (
	"context"
	"log"

	"github.com/gzuuus/onRelay/atomic"
	"github.com/gzuuus/onRelay/policies"
	"github.com/nbd-wtf/go-nostr"
	"github.com/pippellia-btc/rely"
)

/*
Example: Using CombineEvents and CombineFilters with ToClientAdapter for rely
This demonstrates how to combine multiple event policies and use them as a single reject handler,
simplifying registration for frameworks like rely.
*/

var buffer = atomic.NewAtomicCircularBuffer(1000)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go rely.HandleSignals(cancel)

	relay := rely.NewRelay()
	relay.OnEvent = Save
	relay.OnFilters = Query

	// Combine multiple event policies
	combinedEventPolicy := policies.CombineEvents(
		policies.RestrictToSpecifiedKinds(1, 30023),
		policies.RejectEventsWithBase64Media,
	)
	// Use ToClientAdapter to adapt to rely signature
	relay.RejectEvent = append(relay.RejectEvent, policies.ToClientAdapter[rely.Client](combinedEventPolicy))

	addr := "localhost:3335"
	log.Printf("running relay on %s", addr)

	if err := relay.StartAndServe(ctx, addr); err != nil {
		panic(err)
	}
}

func Save(c *rely.Client, e *nostr.Event) error {
	log.Printf("received event: %v", e)
	ctx := context.Background()
	buffer.SaveEvent(ctx, e)
	return nil
}

func Query(ctx context.Context, c *rely.Client, filters nostr.Filters) ([]nostr.Event, error) {
	log.Printf("received filters %v", filters)
	result := make([]nostr.Event, 0)

	for _, f := range filters {
		events, err := buffer.QueryEvents(ctx, f)
		if err != nil {
			log.Printf("[ERROR] querying ephemeral events: %v", err)
		} else {
			for _, event := range events {
				if event != nil {
					result = append(result, *event)
				}
			}
		}
	}
	return result, nil
}
