package hub

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkPublishOnNonBlockingSubscribers(b *testing.B) {
	runBenchmark(b, 100, 4, 30, false)
}

func BenchmarkPublishOnBlockingSubscribers(b *testing.B) {
	runBenchmark(b, 100, 4, 30, true)
}

func runBenchmark(b *testing.B, numItems, numThreads, numSubscribers int, blocking bool) {
	h := New()
	subs := createSubscribers(h, numSubscribers, blocking)
	itemsToInsert := generateTopics(numThreads, numItems)
	var wgPub, wgSub sync.WaitGroup

	wgSub.Add(len(subs))
	processSubscriptionsForBench(subs, &wgSub)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wgPub.Add(numThreads)

		for j := 0; j < numThreads; j++ {
			go func(j int) {
				for _, key := range itemsToInsert[j] {
					h.Publish(Message{Name: key})
				}

				wgPub.Done()
			}(j)
		}
		wgPub.Wait()
	}
	h.Close()
	wgSub.Wait()
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
