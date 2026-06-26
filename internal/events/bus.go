package events

import (
	"context"
	"sync"
)

// Bus is a fan-out pub/sub event bus backed by Go channels.
// One goroutine publishes; many goroutines subscribe and receive
// independent copies of each event.
type Bus struct {
	ch          chan Event
	mu          sync.RWMutex
	subscribers []chan Event
}

func NewBus() *Bus {
	return &Bus{
		// Buffered so publishers never block on a slow dispatch loop.
		ch: make(chan Event, 64),
	}
}

// Publish sends an event into the bus. Non-blocking as long as the
// internal buffer (64) is not full.
func (b *Bus) Publish(e Event) {
	b.ch <- e
}

// Subscribe returns a receive-only channel that will deliver every
// event published after this call. Each subscriber gets its own channel,
// so they never compete with each other.
func (b *Bus) Subscribe() <-chan Event {
	ch := make(chan Event, 16)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Run starts the dispatch loop. It must be called in its own goroutine.
// It exits cleanly when ctx is cancelled.
func (b *Bus) Run(ctx context.Context) {
	for {
		select {
		case e := <-b.ch:
			b.mu.RLock()
			for _, sub := range b.subscribers {
				sub <- e
			}
			b.mu.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}
