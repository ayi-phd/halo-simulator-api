package events_test

import (
	"context"
	"testing"
	"time"

	"halo-simulator/internal/events"
)

func TestBus_DeliverToSubscriber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := events.NewBus()
	go bus.Run(ctx)

	sub := bus.Subscribe()
	bus.Publish(events.Event{Type: events.FlightArrived})

	select {
	case got := <-sub:
		if got.Type != events.FlightArrived {
			t.Errorf("want %s, got %s", events.FlightArrived, got.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestBus_FanOut_AllSubscribersReceive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := events.NewBus()
	go bus.Run(ctx)

	sub1 := bus.Subscribe()
	sub2 := bus.Subscribe()
	sub3 := bus.Subscribe()

	bus.Publish(events.Event{Type: events.FlightArrived})

	timeout := time.After(time.Second)
	for i, sub := range []<-chan events.Event{sub1, sub2, sub3} {
		select {
		case e := <-sub:
			if e.Type != events.FlightArrived {
				t.Errorf("subscriber %d: want %s, got %s", i+1, events.FlightArrived, e.Type)
			}
		case <-timeout:
			t.Fatalf("subscriber %d: timed out waiting for event", i+1)
		}
	}
}
