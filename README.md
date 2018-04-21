# Hub

:incoming_envelope: A fast enough Event Hub for go applications using publish/subscribe with support patterns on topics like rabbitMQ exchanges.

[![Release](https://img.shields.io/github/release/leandro-lugaresi/hub.svg?style=flat-square)](https://github.com/leandro-lugaresi/hub/releases/latest)
[![Software License](https://img.shields.io/github/license/leandro-lugaresi/hub.svg?style=flat-square)](LICENSE.md)
[![Build Status](https://travis-ci.org/leandro-lugaresi/hub.svg?branch=master&style=flat-square)](https://travis-ci.org/leandro-lugaresi/hub)
[![Coverage Status](https://img.shields.io/codecov/c/github/leandro-lugaresi/hub/master.svg?style=flat-square)](https://codecov.io/gh/leandro-lugaresi/hub)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/leandro-lugaresi/hub)
[![Go Report Card](https://goreportcard.com/badge/github.com/leandro-lugaresi/hub?style=flat-square)](https://goreportcard.com/report/github.com/leandro-lugaresi/hub)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/leandro-lugaresi)

---

## Table of Contents

-   [Install](#install)
-   [Usage](#usage)
-   [Examples & Demos](#examples--demos)
-   [Contribute](CONTRIBUTING.md)
-   [Code of conduct](CODE_OF_CONDUCT.md)

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
	"time"

	"github.com/leandro-lugaresi/hub"
)
func main() {
	h := hub.New()

		// the cap param is used to create one buffered channel with cap = 10
		// If you wan an unbuferred channel use the 0 cap
		sub := h.Subscribe("account.*.failed", 10)
		go func(s *hub.Subscription) {
			for msg := range s.Receiver {
				fmt.Printf("receive msg with topic %s and id %d\n", msg.Name, msg.Int("id"))
			}
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
		// TODO: remove sleep when the close method is implemented
		time.Sleep(time.Second)
		// Output:
		// receive msg with topic account.login.failed and id 123
		// receive msg with topic account.changepassword.failed and id 456
		// receive msg with topic account.foo.failed and id 789
}
```

---

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md) for further information
