# Hub

:incoming_envelope: A fast(?) Event Hub using publish/subscribe pattern with support for topics like rabbitMQ exchanges for Go applications

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

This project uses simple go code and you can `go get` it but I encourage you to always vendor your dependencies or use one of the version tags of this project.

```sh
go get -u github.com/leandro-lugaresi/hub
```
```sh
dep ensure --add github.com/leandro-lugaresi/hub
```

## Usage

Hub provides two types of subscribers:
- `BlockingSubscriber` this subscriber use channels inside. If the channel is full the operation of publishing will block.
    you can use a buffered (cap > `0`) or unbuffered channel (cap = `0`).
- `NonBlockingSubscriber` this subscriber use a ring buffer. If the cap of the buffer is reached the publish operation will override the oldest data and never block.
    This should be used only if loose data is acceptable. ie: metrics, logs

## Examples & Demos

TODO: Add real examples

```go

import github.com/leandro-lugaresi/hub

h := hub.New()

// Unbuffered subscribe
subs1 := h.Subscribe("account.login.*", 0)
go func(s *Subscription) {
	for {
		msg, ok := s.Next()
		if !ok {
			return
		}
		//process msg
	}
}(subs1)

subs2 := h.NonBlockingSubscribe("account.login.*", 10)
go func(s *Subscription) {
	for {
		msg, ok := s.Next()
		if !ok {
			return
		}
		//process msg
			}
}(subs2)

h.Publish(hub.Message{
	Name: "account.login.failed",
	Fields: hub.Fields{
		"account": 123,
		"ip": '127.0.0.1',
	},
})
```

---

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md) for further information
