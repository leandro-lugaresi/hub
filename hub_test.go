package hub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	h := New()
	subs1 := h.Subscribe("a.*.c", 0)
	processed := false
	go processSubscription(subs1, func(msg Message) {
		require.Equal(t, "a.b.c", msg.Name)
		require.Equal(t, []byte(`{"foo": "baz"}`), msg.Msg)
		processed = true
	})
	h.Publish(Message{
		Msg:  []byte(`{"foo": "baz"}`),
		Name: "a.b.c",
	})
	time.Sleep(time.Millisecond)
	require.True(t, processed, "subscription function should be executed")

}

func TestNonBlockingSubscribe(t *testing.T) {
	h := New()
	subs1 := h.NonBlockingSubscribe("a.*.c", 10)
	processed := false
	go processSubscription(subs1, func(msg Message) {
		require.Equal(t, "a.b.c", msg.Name)
		require.Equal(t, []byte(`{"foo": "baz"}`), msg.Msg)
		processed = true
	})
	h.Publish(Message{
		Msg:  []byte(`{"foo": "baz"}`),
		Name: "a.b.c",
	})
	time.Sleep(30 * time.Millisecond)
	require.True(t, processed, "subscription function should be executed")

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
