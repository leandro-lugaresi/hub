package hub

const AlertTopic = "hub.subscription.messageslost"

type (
	//Hub is a component that provides publish and subscribe capabilities for messages.
	// Every message has a Name used to route them to subscribers and this can be used like RabbitMQ topics exchanges.
	// Where every word is separated by dots `.` and you can use `*` as a wildcard.
	Hub struct {
		matcher matcher
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
		sub.Set(m)
	}
}

// Subscribe create a blocking subscription to receive events for a given topic.
// The cap param is used inside the subscriber and in this case used to create a channel.
// cap(1) = unbuffered channel.
func (h *Hub) Subscribe(topic string, cap int) Subscription {
	return h.matcher.Subscribe(topic, newBlockingSubscriber(cap))
}

// NonBlockingSubscribe create a nonblocking subscription to receive events for a given topic.
// This implementation use internally a ring buffer so if the buffer reaches the max capability the subscribe will override old messages.
func (h *Hub) NonBlockingSubscribe(topic string, cap int) Subscription {
	return h.matcher.Subscribe(
		topic,
		newNonBlockingSubscriber(
			cap,
			AlertFunc(func(missed int) {
				h.alert(missed, topic)
			}),
		))
}

// Unsubscribe remove and close the Subscription.
func (h *Hub) Unsubscribe(sub Subscription) {
	h.matcher.Unsubscribe(sub)
}

// Close will unsubscribe all the subscriptions and close them all.
func (h *Hub) Close() {
	subs := h.matcher.Subscriptions()
	for _, s := range subs {
		h.matcher.Unsubscribe(s)
	}
	for _, s := range subs {
		s.subscriber.Close()
	}
}

func (h *Hub) alert(missed int, topic string) {
	h.Publish(Message{
		Name: AlertTopic,
		Fields: Fields{
			"missed": missed,
			"topic":  topic,
		},
	})
}
