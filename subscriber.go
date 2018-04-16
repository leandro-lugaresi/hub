package hub

import (
	gendiodes "code.cloudfoundry.org/go-diodes"
)

type (
	// nonBlockingSubscriber uses an diode and is optiomal for many writes and a single reader
	// This subscriber is used when need high throughput and losing data is acceptable.
	nonBlockingSubscriber struct {
		d *gendiodes.Poller
	}

	// blockingSubscriber uses an channel to receive events.
	blockingSubscriber struct {
		ch chan Message
	}

	discardSubscriber int
)

// NewNonBlockingSubscriber returns a new NonBlockingSubscriber diode to be used
// with many writers and a single reader.
func NewNonBlockingSubscriber(cap int, alerter gendiodes.Alerter) Subscriber {
	if cap <= 0 {
		cap = 10
	}
	return &nonBlockingSubscriber{
		d: gendiodes.NewPoller(
			gendiodes.NewManyToOne(cap, alerter)),
	}
}

// Set inserts the given Event into the diode.
func (d *nonBlockingSubscriber) Set(data Message) {
	d.d.Set(gendiodes.GenericDataType(&data))
}

// Next will return the next Event. If the
// diode is empty this method will block until a Event is available to be
// read or context is done. In case of context done we will return false on the second return param.
func (d *nonBlockingSubscriber) Next() (Message, bool) {
	data := d.d.Next()
	if data == nil {
		return Message{}, false
	}
	return *(*Message)(data), true
}

// NewBlockingSubscriber returns a new blocking subscriber using chanels imternally.
func NewBlockingSubscriber(cap int) Subscriber {
	return &blockingSubscriber{
		ch: make(chan Message, cap),
	}
}

// Set inserts the given Event into the diode.
func (s *blockingSubscriber) Set(msg Message) {
	s.ch <- msg
}

// Next will return the next Event. If the
// diode is empty this method will block until a Event is available to be
// read or context is done. In case of context done we will return false on the second return param.
func (s *blockingSubscriber) Next() (Message, bool) {
	msg, ok := <-s.ch
	return msg, ok
}

func (d discardSubscriber) Set(msg Message)       {}
func (d discardSubscriber) Next() (Message, bool) { return Message{}, false }
