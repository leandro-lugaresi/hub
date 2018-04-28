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

func TestHub(t *testing.T) {
	h := New()

	sub0 := h.Subscribe(0, "forex.*")
	sub1 := h.Subscribe(10, "*.usd")
	sub2 := h.Subscribe(-10, "forex.eur", "forex.*")
	sub3 := h.NonBlockingSubscribe(0, "*.eur", "trade")
	sub4 := h.NonBlockingSubscribe(10, "forex.*")
	sub5 := h.NonBlockingSubscribe(-10, "trade")
	sub6 := h.Subscribe(10, "*")

	c0 := newMessageCounter(sub0)
	c1 := newMessageCounter(sub1)
	c2 := newMessageCounter(sub2)
	c3 := newMessageCounter(sub3)
	c4 := newMessageCounter(sub4)
	c5 := newMessageCounter(sub5)
	c6 := newMessageCounter(sub6)

	h.Publish(Message{Name: "forex.eur"})
	h.Publish(Message{Name: "forex"})
	h.Publish(Message{Name: "trade.jpy"})
	h.Publish(Message{Name: "forex.jpy"})
	h.Publish(Message{Name: "trade"})

	time.Sleep(time.Millisecond)

	require.Equal(t, int64(2), c0.count(), "Messages processed by sub0")
	require.Equal(t, int64(0), c1.count(), "Messages processed by sub1")
	require.Equal(t, int64(2), c2.count(), "Messages processed by sub2")
	require.Equal(t, int64(2), c3.count(), "Messages processed by sub3")
	require.Equal(t, int64(2), c4.count(), "Messages processed by sub4")
	require.Equal(t, int64(1), c5.count(), "Messages processed by sub5")
	require.Equal(t, int64(2), c6.count(), "Messages processed by sub6")

	c0.reset()
	c1.reset()
	c2.reset()
	c3.reset()
	c4.reset()
	c5.reset()
	c6.reset()

	h.Close()

	h.Publish(Message{Name: "forex.eur"})
	h.Publish(Message{Name: "forex"})
	h.Publish(Message{Name: "trade.jpy"})
	h.Publish(Message{Name: "forex.jpy"})
	h.Publish(Message{Name: "trade"})

	require.Equal(t, int64(0), c0.count(), "Messages processed by sub0")
	require.Equal(t, int64(0), c1.count(), "Messages processed by sub1")
	require.Equal(t, int64(0), c2.count(), "Messages processed by sub2")
	require.Equal(t, int64(0), c3.count(), "Messages processed by sub3")
	require.Equal(t, int64(0), c4.count(), "Messages processed by sub4")
	require.Equal(t, int64(0), c5.count(), "Messages processed by sub5")
	require.Equal(t, int64(0), c6.count(), "Messages processed by sub6")
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
