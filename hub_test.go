package hub

import (
	"sync"
	"testing"
)

func TestProcessSubscribers(t *testing.T) {
	tests := []struct {
		name   string
		cap   int
		blocking bool
	}{
		{name: "blocking and unbuffered", cap: 0, blocking: true},
		{name: "blocking and buffered", cap: 10, blocking: true},
		{name: "blocking and negative buffer", cap: -10, blocking: true},
		{name: "nonBlocking and unbuffered", cap: 0, blocking: false},
		{name: "nonBlocking and buffered", cap: 10, blocking: false},
		{name: "nonBlocking and negative buffer", cap: -10, blocking: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			var subs *Subscription
			if tt.blocking {
				subs = h.Subscribe("a.*.c", tt.cap)
			} else {
				subs = h.NonBlockingSubscribe("a.*.c", tt.cap)
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go processSubscription(subs, func(msg Message) {
				wg.Done()
			})
			h.Publish(Message{Name: "a.b.c"})
			h.Publish(Message{Name: "a.c.c"})
			wg.Wait()
		})
	}
}

func processSubscription(s *Subscription, op func(msg Message)) {
	for {
		msg, ok := s.Subscriber.Next()
		if !ok {
			break
		}
		op(msg)
	}
}
