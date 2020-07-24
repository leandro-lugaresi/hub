# Hub

:incoming_envelope: A fast enough Event Hub for go applications using publish/subscribe with support patterns on topics like rabbitMQ exchanges.

[![Release](https://img.shields.io/github/release/leandro-lugaresi/hub.svg?style=flat-square)](https://github.com/leandro-lugaresi/hub/releases/latest)
[![Software License](https://img.shields.io/github/license/leandro-lugaresi/hub.svg?style=flat-square)](LICENSE.md)
[![Actions Status](https://github.com/leandro-lugaresi/hub/workflows/Go/badge.svg)](https://github.com/leandro-lugaresi/hub/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/leandro-lugaresi/hub/main.svg?style=flat-square)](https://codecov.io/gh/leandro-lugaresi/hub)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/leandro-lugaresi/hub)
[![Go Report Card](https://goreportcard.com/badge/github.com/leandro-lugaresi/hub?style=flat-square)](https://goreportcard.com/report/github.com/leandro-lugaresi/hub)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/leandro-lugaresi)

---

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Examples & Demos](#examples--demos)
- [Benchmarks](#benchmarks)
- [CSTrie](#cstrie)
- [Contribute](CONTRIBUTING.md)
- [Code of conduct](CODE_OF_CONDUCT.md)

## Install

To install this library you can `go get` it but I encourage you to always vendor your dependencies or use one of the version tags of this project.

```sh
go get -u github.com/leandro-lugaresi/hub
```

```sh
dep ensure --add github.com/leandro-lugaresi/hub
```

## Usage

### Subscribers

Hub provides subscribers as buffered (cap > `0`) and unbuffered (cap = 0) channels but we have two different types of subscribers:

- `Subscriber` this is the default subscriber and it's a blocking subscriber so if the channel is full and you try to send another message the send operation will block until the subscriber consumes some message.
- `NonBlockingSubscriber` this subscriber will never block on the publish side but if the capacity of the channel is reached the publish operation will be lost and an alert will be trigged.
  This should be used only if loose data is acceptable. ie: metrics, logs

### Topics

This library uses the same concept of topic exchanges on rabbiMQ, so the message name is used to find all the subscribers that match the topic, like a route.
The topic must be a list of words delimited by dots (`.`) however, there is one important special case for binding keys:
`*` (star) can substitute for exactly one word.

## Examples & Demos

```go
package main

import (
	"fmt"
	"sync"

	"github.com/leandro-lugaresi/hub"
)
func main() {
	h := hub.New()
	var wg sync.WaitGroup
	// the cap param is used to create one buffered channel with cap = 10
	// If you wan an unbuferred channel use the 0 cap
	sub := h.Subscribe(10, "account.login.*", "account.changepassword.*")
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
		Name:   "account.login.success",
		Fields: hub.Fields{"id": 123},
	})
	// message not routed to this subscriber
	h.Publish(hub.Message{
		Name:   "account.foo.failed",
		Fields: hub.Fields{"id": 789},
	})

	// close all the subscribers
	h.Close()
	// wait until finish all the messages on buffer
	wg.Wait()

	// Output:
	// receive msg with topic account.login.failed and id 123
	// receive msg with topic account.changepassword.failed and id 456
	// receive msg with topic account.login.success and id 123
}
```

See more [here](https://godoc.org/github.com/leandro-lugaresi/hub#example-Hub)!

## Benchmarks

To run the benchmarks you can execute:

```bash
make bench
```

Currently, I only have the benchmarks of the CSTrie used internally. I will provide more benchmarks.

## Throughput

The project have one test for throughput, just execute:

```bash
make throughput
```

In a intel(R) core(TM) i5-4460 CPU @ 3.20GHz x4 we got this results:

```
go test -v -timeout 60s github.com/leandro-lugaresi/hub -run ^TestThroughput -args -throughput
=== RUN   TestThroughput
1317530.091292 msg/sec
--- PASS: TestThroughput (3.04s)
PASS
ok      github.com/leandro-lugaresi/hub 3.192s
```

## CSTrie

This project uses internally an awesome Concurrent Subscription Trie done by [@tylertreat](https://github.com/tylertreat). If you want to learn more about see this [blog post](http://bravenewgeek.com/fast-topic-matching/) and the code is [here](https://github.com/tylertreat/fast-topic-matching)

---

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md) for further information
