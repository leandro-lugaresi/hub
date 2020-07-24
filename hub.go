package hub

// AlertTopic is used to notify when a nonblocking subscriber loose one message
// You can subscribe on this topic and log or send metrics.
const AlertTopic = "hub.subscription.messageslost"

type (
	//Hub is a component that provides publish and subscribe capabilities for messages.
	// Every message has a Name used to route them to subscribers and this can be used like RabbitMQ topics exchanges.
	// Where every word is separated by dots `.` and you can use `*` as a wildcard.
	Hub struct {
		matcher matcher
		fields  Fields
	}
)

// New create and return a new empty hub.
func New() *Hub {
	return &Hub{
		matcher: newCSTrieMatcher(),
		fields:  Fields{},
	}
}

// Publish will send an event to all the subscribers matching the event name.
func (h *Hub) Publish(m Message) {
	for k, v := range h.fields {
		m.Fields[k] = v
	}

	for _, sub := range h.matcher.Lookup(m.Topic()) {
		sub.Set(m)
	}
}

// With creates a child Hub with the fields added to it.
// When someone call Publish, this Fields will be added automatically into the message.
func (h *Hub) With(f Fields) *Hub {
	hub := Hub{
		matcher: h.matcher,
		fields:  Fields{},
	}
	for k, v := range h.fields {
		hub.fields[k] = v
	}

	for k, v := range f {
		hub.fields[k] = v
	}

	return &hub
}

// Subscribe create a blocking subscription to receive events for a given topic.
// The cap param is used inside the subscriber and in this case used to create a channel.
// cap(1) = unbuffered channel.
func (h *Hub) Subscribe(cap int, topics ...string) Subscription {
	return h.matcher.Subscribe(topics, newBlockingSubscriber(cap))
}

// NonBlockingSubscribe create a nonblocking subscription to receive events for a given topic.
// This subscriber will loose messages if the buffer reaches the max capability.
func (h *Hub) NonBlockingSubscribe(cap int, topics ...string) Subscription {
	return h.matcher.Subscribe(
		topics,
		newNonBlockingSubscriber(
			cap,
			alertFunc(func(missed int) {
				h.alert(missed, topics)
			}),
		))
}

// Unsubscribe remove and close the Subscription.
func (h *Hub) Unsubscribe(sub Subscription) {
	h.matcher.Unsubscribe(sub)
	sub.subscriber.Close()
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

func (h *Hub) alert(missed int, topics []string) {
	h.Publish(Message{
		Name: AlertTopic,
		Fields: Fields{
			"missed": missed,
			"topic":  topics,
		},
	})
}
