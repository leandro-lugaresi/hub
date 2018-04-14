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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSTrieMatcher(t *testing.T) {
	assert := assert.New(t)
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
		s1 = discardSubscriber(1)
		s2 = discardSubscriber(2)
	)

	sub0, err := m.Subscribe("forex.*", s0)
	assert.NoError(err)
	sub1, err := m.Subscribe("*.usd", s0)
	assert.NoError(err)
	sub2, err := m.Subscribe("forex.eur", s0)
	assert.NoError(err)
	sub3, err := m.Subscribe("*.eur", s1)
	assert.NoError(err)
	sub4, err := m.Subscribe("forex.*", s1)
	assert.NoError(err)
	sub5, err := m.Subscribe("trade", s1)
	assert.NoError(err)
	sub6, err := m.Subscribe("*", s2)
	assert.NoError(err)

	assertEqual(assert, []Subscriber{s0, s1}, m.Lookup("forex.eur"))
	assertEqual(assert, []Subscriber{s2}, m.Lookup("forex"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Subscriber{s0, s1}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Subscriber{s1, s2}, m.Lookup("trade"))

	m.Unsubscribe(sub0)
	m.Unsubscribe(sub1)
	m.Unsubscribe(sub2)
	m.Unsubscribe(sub3)
	m.Unsubscribe(sub4)
	m.Unsubscribe(sub5)
	m.Unsubscribe(sub6)

	assertEqual(assert, []Subscriber{}, m.Lookup("forex.eur"))
	assertEqual(assert, []Subscriber{}, m.Lookup("forex"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Subscriber{}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade"))
}

func BenchmarkCSTrieMatcherSubscribe(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)
	populateMatcher(m, 1000, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Subscribe("foo.*.baz.qux.quux", s0)
	}
}

func BenchmarkCSTrieMatcherUnsubscribe(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)
	id, _ := m.Subscribe("foo.*.baz.qux.quux", s0)
	populateMatcher(m, 1000, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Unsubscribe(id)
	}
}

func BenchmarkCSTrieMatcherLookup(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)
	m.Subscribe("foo.*.baz.qux.quux", s0)
	populateMatcher(m, 1000, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lookup("foo.bar.baz.qux.quux")
	}
}

func BenchmarkCSTrieMatcherSubscribeCold(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Subscribe("foo.*.baz.qux.quux", s0)
	}
}

func BenchmarkCSTrieMatcherUnsubscribeCold(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)
	id, _ := m.Subscribe("foo.*.baz.qux.quux", s0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Unsubscribe(id)
	}
}

func BenchmarkCSTrieMatcherLookupCold(b *testing.B) {
	var (
		m  = NewCSTrieMatcher()
		s0 = discardSubscriber(0)
	)
	m.Subscribe("foo.*.baz.qux.quux", s0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lookup("foo.bar.baz.qux.quux")
	}
}

func BenchmarkMultithreaded1Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 1
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded2Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 2
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded4Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 4
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded8Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 8
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded12Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 12
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded16Thread5050CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 16
	benchmark5050(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded1Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 1
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded2Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 2
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded4Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 4
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded8Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 8
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded12Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 12
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}

func BenchmarkMultithreaded16Thread9010CSTrie(b *testing.B) {
	numItems := 1000
	numThreads := 16
	benchmark9010(b, numItems, numThreads, func(items [][]string) Matcher {
		return NewCSTrieMatcher()
	})
}
