// Copyright (C) 2018 Tyler Treat <https://github.com/tylertreat>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Modifications copyright (C) 2018 Leandro Lugaresi

package hub

const (
	delimiter = "."
	wildcard  = "*"
)

type (
	// Subscription represents a topic subscription.
	Subscription struct {
		Topics     []string
		Receiver   <-chan Message
		subscriber subscriber
	}

	// subscriber is the interface used internally to send values and get the channel used by subscribers.
	// This is used to override the behaviour of channel and support nonBlocking operations
	subscriber interface {
		// Set send the given Event to be processed by the subscriber
		Set(Message)
		// Ch return the channel used to consume messages inside the subscription.
		// This func MUST always return the same channel.
		Ch() <-chan Message
		// Close will close the internal state and the subscriber will not receive more messages
		// WARN: This function can be executed more than one time so the code MUST take care of this situation and
		// avoid problems like close a closed channel.
		Close()
	}
)

// matcher contains topic subscriptions and performs matches on them.
type matcher interface {
	// Subscribe adds the Subscriber to the topics and returns a Subscription.
	Subscribe(topics []string, sub subscriber) Subscription

	// Unsubscribe removes the Subscription.
	Unsubscribe(sub Subscription)

	// Lookup returns the subscribers for the given topic.
	Lookup(topic string) []subscriber

	Subscriptions() []Subscription
}
