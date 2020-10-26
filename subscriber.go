package hub

import (
	"sync"
)

type (
	alertFunc func(missed int)

	nonBlockingSubscriber struct {
		ch     chan Message
		closed bool
		// chMu protects ch and closed
		chMu      sync.RWMutex
		alert     alertFunc
		onceClose sync.Once
	}
	// blockingSubscriber uses an channel to receive events.
	blockingSubscriber struct {
		ch        chan Message
		chMu      sync.RWMutex
		onceClose sync.Once

		closeCh chan struct{}
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
	s.chMu.RLock()
	defer s.chMu.RUnlock()

	if s.closed {
		s.alert(1)
		return
	}

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
		s.chMu.Lock()
		defer s.chMu.Unlock()

		s.closed = true

		close(s.ch)
	})
}

// newBlockingSubscriber returns a blocking subscriber using chanels imternally.
func newBlockingSubscriber(cap int) *blockingSubscriber {
	if cap < 0 {
		cap = 0
	}

	return &blockingSubscriber{
		ch:      make(chan Message, cap),
		closeCh: make(chan struct{}),
	}
}

// Set will send the message using the channel.
func (s *blockingSubscriber) Set(msg Message) {
	s.chMu.RLock()
	defer s.chMu.RUnlock()

	// check s.ch isn't closed (we are holding the RLock, so s.ch won't be closed until the end of this function)
	select {
	case <-s.closeCh:
		return
	default:
	}

	s.ch <- msg
}

// Ch return the channel used by subscriptions to consume messages.
func (s *blockingSubscriber) Ch() <-chan Message {
	return s.ch
}

// Close will close the internal channel and stop receiving messages.
func (s *blockingSubscriber) Close() {
	s.onceClose.Do(func() {
		// make sure we can take the mutex (nobody is holding a RLock in a blocking Set)
		close(s.closeCh)

		s.chMu.Lock()
		defer s.chMu.Unlock()

		close(s.ch)
	})
}
