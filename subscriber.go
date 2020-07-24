package hub

import (
	"sync"
)

type (
	alertFunc func(missed int)

	nonBlockingSubscriber struct {
		ch        chan Message
		alert     alertFunc
		onceClose sync.Once
	}
	// blockingSubscriber uses an channel to receive events.
	blockingSubscriber struct {
		ch        chan Message
		onceClose sync.Once
	}
)

// newNonBlockingSubscriber returns a new nonBlockingSubscriber
// this subscriber will never block when sending an message, if the capacity is full
// we will ignore the message and call the Alert function from the Alerter.
func newNonBlockingSubscriber(cap int, alerter alertFunc) *nonBlockingSubscriber {
	if cap <= 0 {
		cap = 10
	}

	return &nonBlockingSubscriber{
		ch:    make(chan Message, cap),
		alert: alerter,
	}
}

// Set inserts the given Event into the diode.
func (s *nonBlockingSubscriber) Set(msg Message) {
	select {
	case s.ch <- msg:
	default:
		s.alert(1)
	}
}

// Ch return the channel used by subscriptions to consume messages.
func (s *nonBlockingSubscriber) Ch() <-chan Message {
	return s.ch
}

// Close will close the internal channel and stop receiving messages.
func (s *nonBlockingSubscriber) Close() {
	s.onceClose.Do(func() {
		close(s.ch)
	})
}

// newBlockingSubscriber returns a blocking subscriber using chanels imternally.
func newBlockingSubscriber(cap int) *blockingSubscriber {
	if cap < 0 {
		cap = 0
	}

	return &blockingSubscriber{
		ch: make(chan Message, cap),
	}
}

// Set will send the message using the channel.
func (s *blockingSubscriber) Set(msg Message) {
	s.ch <- msg
}

// Ch return the channel used by subscriptions to consume messages.
func (s *blockingSubscriber) Ch() <-chan Message {
	return s.ch
}

// Close will close the internal channel and stop receiving messages.
func (s *blockingSubscriber) Close() {
	s.onceClose.Do(func() {
		close(s.ch)
	})
}
