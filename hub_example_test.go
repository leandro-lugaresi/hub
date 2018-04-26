package hub_test

import (
	"fmt"
	"sync"

	"github.com/leandro-lugaresi/hub"
)

func ExampleHub() {
	h := hub.New()
	var wg sync.WaitGroup
	// the cap param is used to create one buffered channel with cap = 10
	// If you wan an unbuferred channel use the 0 cap
	sub := h.Subscribe("account.*.failed", 10)
	wg.Add(1)
	go func(s hub.Subscription) {
		for msg := range s.Receiver {
			fmt.Printf("receive msg with topic %s and id %d\n", msg.Name, msg.Fields["id"])
		}
		wg.Done()
	}(sub)

	h.Publish(hub.Message{
		Name:   "account.login.failed",
		Fields: hub.Fields{"id": 123},
	})
	h.Publish(hub.Message{
		Name:   "account.changepassword.failed",
		Fields: hub.Fields{"id": 456},
	})
	h.Publish(hub.Message{
		Name:   "account.foo.failed",
		Fields: hub.Fields{"id": 789},
	})

	// finish all the subscribers
	h.Close()
	// wait until finish all the messages on buffer
	wg.Wait()

	// Output:
	// receive msg with topic account.login.failed and id 123
	// receive msg with topic account.changepassword.failed and id 456
	// receive msg with topic account.foo.failed and id 789
}
