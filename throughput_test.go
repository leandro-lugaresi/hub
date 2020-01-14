package hub

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	numSubs       = 1000
	numMsgs       = 1000000
	numPublishers = 4
)

var (
	throughputTest = flag.Bool("throughput", false, "execute throughput tests")

	topics = make([]string, numSubs)
	msgs   = make([]Message, numMsgs)
)

func setupTopics() {
	for i := 0; i < numSubs; i++ {
		switch {
		case i%10 == 0:
			topics[i] = fmt.Sprintf("*.%d.%d", rand.Intn(10), rand.Intn(10))
		case i%25 == 0:
			topics[i] = fmt.Sprintf("%d.*.%d", rand.Intn(10), rand.Intn(10))
		case i%45 == 0:
			topics[i] = fmt.Sprintf("%d.%d.*", rand.Intn(10), rand.Intn(10))
		default:
			topics[i] = fmt.Sprintf("%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10))
		}
	}

	for i := 0; i < numMsgs; i++ {
		topic := topics[i%numSubs]
		msgs[i] = Message{
			Name: strings.Replace(topic, "*", strconv.Itoa(rand.Intn(10)), -1),
		}
	}
}

func TestThroughput(t *testing.T) {
	setupTopics()

	h := New()
	var wg sync.WaitGroup

	if !*throughputTest {
		t.Skip("throughputTests skipped")
	}

	for _, topic := range topics {
		sub := h.Subscribe(200, topic)
		go func(s Subscription) {
			for range s.Receiver {
			}

			wg.Done()
		}(sub)
	}

	before := time.Now()

	wg.Add(numPublishers)

	for i := 0; i < numPublishers; i++ {
		go func() {
			for _, msg := range msgs {
				h.Publish(msg)
			}

			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(numSubs)
	h.Close()
	wg.Wait()

	dur := time.Since(before)
	throughput := numMsgs * numPublishers / dur.Seconds()
	fmt.Printf("%f msg/sec\n", throughput)
}
