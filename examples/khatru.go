package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fiatjaf/khatru"
	"github.com/gzuuus/onRelay/atomic"
	"github.com/nbd-wtf/go-nostr"
)

var buffer = atomic.NewAtomicCircularBuffer(1000)

func main() {
	relay := khatru.NewRelay()

	relay.StoreEvent = append(relay.StoreEvent, buffer.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, atomic.SliceToChanEvents(buffer.QueryEvents))

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
