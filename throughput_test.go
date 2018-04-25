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

func init() {
	for i := 0; i < numSubs; i++ {
		if i%10 == 0 {
			topics[i] = fmt.Sprintf("*.%d.%d", rand.Intn(10), rand.Intn(10))
		} else if i%25 == 0 {
			topics[i] = fmt.Sprintf("%d.*.%d", rand.Intn(10), rand.Intn(10))
		} else if i%45 == 0 {
			topics[i] = fmt.Sprintf("%d.%d.*", rand.Intn(10), rand.Intn(10))
		} else {
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
	h := New()
	if !*throughputTest {
		t.Skip("throughputTests skipped")
	}
	for _, topic := range topics {
		sub := h.NonBlockingSubscribe(topic, 200)
		go func(s Subscription) {
			for range s.Receiver {

			}
		}(sub)
	}

	before := time.Now()
	var wg sync.WaitGroup
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
	dur := time.Since(before)
	throughput := numMsgs * numPublishers / dur.Seconds()
	fmt.Printf("%f msg/sec\n", throughput)
}
