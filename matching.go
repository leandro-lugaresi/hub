//Copyright (C) 2018 Tyler Treat <https://github.com/tylertreat>
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at

//	http://www.apache.org/licenses/LICENSE-2.0

//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

//Modifications copyright (C) 2018 Leandro Lugaresi

package hub

const (
	delimiter = "."
	wildcard  = "*"
)

// Subscription represents a topic subscription.
type Subscription struct {
	id         uint32
	Topic      string
	Subscriber Subscriber
}

// Matcher contains topic subscriptions and performs matches on them.
type Matcher interface {
	// Subscribe adds the Subscriber to the topic and returns a Subscription.
	Subscribe(topic string, sub Subscriber) *Subscription

	// Unsubscribe removes the Subscription.
	Unsubscribe(sub *Subscription)

	// Lookup returns the SubscrMultithreaded4Threadibers for the given topic.
	Lookup(topic string) []Subscriber
}
