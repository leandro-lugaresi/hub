package hub

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type messageCounter struct {
	c   int64
	sub Subscription
}

// nolint:funlen
func TestHub(t *testing.T) {
	h := New()
	defaultMessages := []Message{
		{Name: "forex.eur"},
		{Name: "forex"},
		{Name: "trade.jpy"},
		{Name: "forex.jpy"},
		{Name: "trade"},
		{Name: "trade.usd"},
		{Name: "forex.usd"},
		{Name: "trade.eur"},
	}
	testCases := []struct {
		name          string
		messages      []Message
		subFN         func(h *Hub) Subscription
		ExpectedCount int
	}{
		{
			name:          "simple subscription unbuffered",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.Subscribe(0, "forex.*") },
			ExpectedCount: 3,
		},
		{
			name:          "simple subscription buffered",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.Subscribe(2, "*.usd") },
			ExpectedCount: 2,
		},
		{
			name:          "simple subscription with an invalid cap buffer",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.Subscribe(-1, "forex", "forex.eur", "forex.*") },
			ExpectedCount: 4,
		},
		{
			name:          "non blocking subscription unbuffered",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.NonBlockingSubscribe(0, "*.eur", "trade") },
			ExpectedCount: 3,
		},
		{
			name:          "non blocking subscription buffered",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.NonBlockingSubscribe(2, "forex", "forex.eur", "forex.*") },
			ExpectedCount: 4,
		},
		{
			name:          "non blocking subscription with an invalid cap buffer",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.NonBlockingSubscribe(-1, "trade") },
			ExpectedCount: 1,
		},
		{
			name:          "get all the messages",
			messages:      defaultMessages,
			subFN:         func(h *Hub) Subscription { return h.NonBlockingSubscribe(0, "*", "*.*") },
			ExpectedCount: 8,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sub := tc.subFN(h)
			counter := newMessageCounter(sub)
			for _, m := range tc.messages {
				time.Sleep(time.Millisecond)
				h.Publish(m)
			}

			time.Sleep(time.Second)
			require.EqualValues(t, tc.ExpectedCount, counter.count())

			counter.reset()
			h.Close()

			for _, m := range tc.messages {
				h.Publish(m)
			}
			require.EqualValues(t, 0, counter.count(), "after close the hub all the subsctibers MUST not get more events")
		})
	}
}

func TestNonBlockingSubscriberShouldAlertIfLoseMessages(t *testing.T) {
	h := New()
	h.NonBlockingSubscribe(10, "a.*.c")
	subsAlert := h.Subscribe(1, AlertTopic)
	// send messages without a working subscriber
	for i := 0; i < 11; i++ {
		h.Publish(Message{Name: "a.c.c", Fields: Fields{"i": i}})
	}

	msg := <-subsAlert.Receiver
	require.Equal(t, 1, msg.Fields["missed"])
	require.Equal(t, []string{"a.*.c"}, msg.Fields["topic"])
}

func TestWith(t *testing.T) {
	h := New()
	subH1 := h.With(Fields{"hub": "subH1", "something": 123})
	subH11 := subH1.With(Fields{"hub": "subH11", "field": 456})
	subH2 := h.With(Fields{"hub": "subH2", "something": 789})

	subs := h.Subscribe(5, "*")

	h.Publish(Message{Name: "foo", Fields: Fields{"msg": 1}})
	subH1.Publish(Message{Name: "foo", Fields: Fields{"msg": 2}})
	subH11.Publish(Message{Name: "foo", Fields: Fields{"msg": 3}})
	subH2.Publish(Message{Name: "foo", Fields: Fields{"msg": 4, "something": 1234}})

	msg := <-subs.Receiver
	require.Equal(t, Fields{"msg": 1}, msg.Fields)
	msg = <-subs.Receiver
	require.Equal(t, Fields{"msg": 2, "hub": "subH1", "something": 123}, msg.Fields)
	msg = <-subs.Receiver
	require.Equal(t, Fields{
		"msg":       3,
		"hub":       "subH11",
		"something": 123,
		"field":     456,
	}, msg.Fields)

	msg = <-subs.Receiver
	require.Equal(t, Fields{"msg": 4, "hub": "subH2", "something": 789}, msg.Fields)
}

func newMessageCounter(s Subscription) *messageCounter {
	ms := &messageCounter{sub: s, c: 0}
	go func(ms *messageCounter) {
		for range ms.sub.Receiver {
			atomic.AddInt64(&ms.c, 1)
		}
	}(ms)

	return ms
}

func (ms *messageCounter) count() int64 {
	return atomic.LoadInt64(&ms.c)
}

func (ms *messageCounter) reset() {
	atomic.StoreInt64(&ms.c, int64(0))
}
