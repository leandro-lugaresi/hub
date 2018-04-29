package hub

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkPublishOnNonBlockingSubscribers(b *testing.B) {
	var wg sync.WaitGroup
	h := New()
	subs := createSubscribers(h, 30, false)
	wg.Add(len(subs))
	processSubscriptionsForBench(subs, &wg)
	runBenchmark(b, 100, 4, h)
	h.Close()
	wg.Wait()
}

func BenchmarkPublishOnBlockingSubscribers(b *testing.B) {
	var wg sync.WaitGroup
	h := New()
	subs := createSubscribers(h, 30, true)
	wg.Add(len(subs))
	processSubscriptionsForBench(subs, &wg)
	runBenchmark(b, 100, 4, h)
	h.Close()
	wg.Wait()
}

func runBenchmark(b *testing.B, numItems, numThreads int, h *Hub) {
	itemsToInsert := generateTopics(numThreads, numItems)

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(numThreads)
		for j := 0; j < numThreads; j++ {
			go func(j int) {
				for _, key := range itemsToInsert[j] {
					h.Publish(Message{Name: key})
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
	}
	b.StopTimer()
}

func createSubscribers(h *Hub, qtd int, blocking bool) []Subscription {
	subs := make([]Subscription, 0, qtd)
	for i := 0; i < qtd; i++ {
		topic := strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(50)) + ".*"
		var sub Subscription
		if blocking {
			sub = h.Subscribe(20, topic)
		} else {
			sub = h.NonBlockingSubscribe(20, topic)
		}
		subs = append(subs, sub)
	}
	return subs
}

func processSubscriptionsForBench(subs []Subscription, wg *sync.WaitGroup) {
	for _, sub := range subs {
		go func(s Subscription) {
			for range s.Receiver {

			}
			wg.Done()
		}(sub)
	}
}
