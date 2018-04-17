package hub

import (
	diodes "code.cloudfoundry.org/go-diodes"
)

type (
	//Hub is a component that provides publish and subscribe capabilities for messages.
	// Every message has a Name used to route them to subscribers and this can be used like RabbitMQ topics exchanges.
	// Where every word is separated by dots `.` and you can use `*` as a wildcard.
	Hub struct {
		matcher Matcher
	}

	// publisher is the interface used internally to send values to the Receiver
	publisher interface {
		// Set send the given Event to be processed by the subscriber
		Set(Message)
	}

	// Receiver is the interface used to get Messages from hub.
	Receiver interface {
		// Next will return the next Event to be processed.
		// This method will block until an Event is available or context is done.
		// In case of context done we will return true on the second return param.
		Next() (Message, bool)
	}

	// Subscriber is the interface used to get Messages from Hub.
	Subscriber interface {
		Receiver
		publisher
	}
)

// New create and return a new empty hub.
func New() *Hub {
	return &Hub{
		matcher: newCSTrieMatcher(),
	}
}

// Publish will send an event to all the subscribers matching the event name.
func (h *Hub) Publish(m Message) {
	for _, sub := range h.matcher.Lookup(m.Topic()) {
		s := sub.(publisher)
		s.Set(m)
	}
}

// Subscribe create a blocking subscription to receive events for a given topic.
// The cap param is used inside the subscriber and in this case used to create a channel.
// cap(1) = unbuffered channel.
func (h *Hub) Subscribe(topic string, cap int) *Subscription {
	return h.matcher.Subscribe(topic, NewBlockingSubscriber(cap))
}

// NonBlockingSubscribe create a nonblocking subscription to receive events for a given topic.
// This implementation use internally a ring buffer so if the buffer reaches the max capability the subscribe will override old messages.
func (h *Hub) NonBlockingSubscribe(topic string, cap int) *Subscription {
	return h.matcher.Subscribe(
		topic,
		NewNonBlockingSubscriber(
			cap,
			diodes.AlertFunc(func(missed int) {
				h.alert(missed, topic)
			}),
		))
}

// Unsubscribe remove and close the Subscription.
func (h *Hub) Unsubscribe(sub *Subscription) {
	h.matcher.Unsubscribe(sub)
}

// Close will remove and unsubcribe all the subscriptions.
// TODO: Implement to close every subscriber nicelly.
func (*Hub) Close() error {
	return nil
}

func (h *Hub) alert(missed int, topic string) {
	h.Publish(Message{
		Name: "hub.subscription.messageslost",
		Fields: Fields{
			"missed": missed,
			"topic":  topic,
		},
	})
}
