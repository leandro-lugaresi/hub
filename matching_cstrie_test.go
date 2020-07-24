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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSTrieMatcher(t *testing.T) {
	assert := assert.New(t)
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
		s1 = discardSubscriber(1)
		s2 = discardSubscriber(2)
	)

	sub0 := m.Subscribe([]string{"forex.*"}, s0)
	sub1 := m.Subscribe([]string{"*.usd"}, s0)
	sub2 := m.Subscribe([]string{"forex.eur"}, s0)
	sub3 := m.Subscribe([]string{"*.eur"}, s1)
	sub4 := m.Subscribe([]string{"forex.*"}, s1)
	sub5 := m.Subscribe([]string{"trade"}, s1)
	sub6 := m.Subscribe([]string{"*"}, s2)
	assert.Len(m.Subscriptions(), 7)

	assertEqual(assert, []subscriber{s0, s1}, m.Lookup("forex.eur"))
	assertEqual(assert, []subscriber{s2}, m.Lookup("forex"))
	assertEqual(assert, []subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []subscriber{s0, s1}, m.Lookup("forex.jpy"))
	assertEqual(assert, []subscriber{s1, s2}, m.Lookup("trade"))

	m.Unsubscribe(sub0)
	m.Unsubscribe(sub1)
	m.Unsubscribe(sub2)
	m.Unsubscribe(sub3)
	m.Unsubscribe(sub4)
	m.Unsubscribe(sub5)
	m.Unsubscribe(sub6)

	assertEqual(assert, []subscriber{}, m.Lookup("forex.eur"))
	assertEqual(assert, []subscriber{}, m.Lookup("forex"))
	assertEqual(assert, []subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []subscriber{}, m.Lookup("forex.jpy"))
	assertEqual(assert, []subscriber{}, m.Lookup("trade"))
}

func BenchmarkCSTrieMatcherSubscribe(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
	)

	populateMatcher(m, 5)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	}
}

func BenchmarkCSTrieMatcherUnsubscribe(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
		id = m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	)

	populateMatcher(m, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Unsubscribe(id)
	}
}

func BenchmarkCSTrieMatcherLookup(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
	)

	m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	populateMatcher(m, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Lookup("foo.bar.baz.qux.quux")
	}
}

func BenchmarkCSTrieMatcherSubscribeCold(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	}
}

func BenchmarkCSTrieMatcherUnsubscribeCold(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
		id = m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Unsubscribe(id)
	}
}

func BenchmarkCSTrieMatcherLookupCold(b *testing.B) {
	var (
		m  = newCSTrieMatcher()
		s0 = discardSubscriber(0)
	)

	m.Subscribe([]string{"foo.*.baz.qux.quux"}, s0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Lookup("foo.bar.baz.qux.quux")
	}
}

func BenchmarkMultithreaded1Thread5050CSTrie(b *testing.B) {
	numThreads := 1
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded2Thread5050CSTrie(b *testing.B) {
	numThreads := 2
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded4Thread5050CSTrie(b *testing.B) {
	numThreads := 4
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded8Thread5050CSTrie(b *testing.B) {
	numThreads := 8
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded12Thread5050CSTrie(b *testing.B) {
	numThreads := 12
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded16Thread5050CSTrie(b *testing.B) {
	numThreads := 16
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual5050)
}

func BenchmarkMultithreaded1Thread9010CSTrie(b *testing.B) {
	numThreads := 1
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}

func BenchmarkMultithreaded2Thread9010CSTrie(b *testing.B) {
	numThreads := 2
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}

func BenchmarkMultithreaded4Thread9010CSTrie(b *testing.B) {
	numThreads := 4
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}

func BenchmarkMultithreaded8Thread9010CSTrie(b *testing.B) {
	numThreads := 8
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}

func BenchmarkMultithreaded12Thread9010CSTrie(b *testing.B) {
	numThreads := 12
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}

func BenchmarkMultithreaded16Thread9010CSTrie(b *testing.B) {
	numThreads := 16
	benchmarkMatcher(b, numThreads, newCSTrieMatcher(), percentual9010)
}
