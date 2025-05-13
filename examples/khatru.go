package main

import (
	"fmt"
	"net/http"

	"github.com/fiatjaf/khatru"
	"github.com/gzuuus/onRelay/atomic"
)

// Create a buffer with capacity for 1000 events
var buffer = atomic.NewAtomicCircularBuffer(1000)

func main() {
	// Initialize a new khatru relay
	relay := khatru.NewRelay()

	// Register the event storage handler
	relay.StoreEvent = append(relay.StoreEvent, buffer.SaveEvent)
	
	// Use the QueryEventsToChan adapter to convert our slice-based QueryEvents
	// to the channel-based signature that khatru expects
	relay.QueryEvents = append(relay.QueryEvents, atomic.QueryEventsToChan(buffer.QueryEvents))

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
