package hub

type (
	AlertFunc func(missed int)

	nonBlockingSubscriber struct {
		ch    chan Message
		alert AlertFunc
	}
	// blockingSubscriber uses an channel to receive events.
	blockingSubscriber chan Message
)

// newNonBlockingSubscriber returns a new nonBlockingSubscriber
// this subscriber will never block when sending an message, if the capacity is full
// we will ignore the message and call the Alert function from the Alerter.
func newNonBlockingSubscriber(cap int, alerter AlertFunc) *nonBlockingSubscriber {
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

// newBlockingSubscriber returns a blocking subscriber using chanels imternally.
func newBlockingSubscriber(cap int) subscriber {
	if cap < 0 {
		cap = 0
	}
	return make(blockingSubscriber, cap)
}

// Set will send the message using the channel
func (ch blockingSubscriber) Set(msg Message) {
	ch <- msg
}

// Ch return the channel used by subscriptions to consume messages
func (ch blockingSubscriber) Ch() <-chan Message {
	return ch
}
