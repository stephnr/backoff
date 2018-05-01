# Backoff
![Travis](https://travis-ci.org/defaltd/backoff.svg?branch=master) [![Coverage Status](https://coveralls.io/repos/github/defaltd/backoff/badge.svg?branch=master)](https://coveralls.io/github/defaltd/backoff?branch=master) [![Godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/defaltd/backoff) [![Go Report Card](https://goreportcard.com/badge/github.com/defaltd/backoff)](https://goreportcard.com/report/github.com/defaltd/backoff)

A package for Optimistic Concurrency Control

## Summary

This project is an improvement of the backoff package implemented by cenkalti/backoff. This package provides support for various other backoff algorithms and exposes them in a service-like fashion for use.

Backoff Algorithms such as Exponential Backoff are all part of a branch of computer science known as [“Optimistic Concurrency Control”](https://en.wikipedia.org/wiki/Optimistic_concurrency_control). The focus of these algorithms is to use feedback cycles and process staggering to determine an acceptable time at which a process can successfully complete its intended goal.

Supported Algorithms:

- [x] [Exponential Backoff](https://en.wikipedia.org/wiki/Exponential_backoff)
- [ ] Fibonacci (TBD)
- [ ] Equal Jitter (TBD)
- [ ] Full Jitter (TBD)
- [ ] Decorr (TBD)
- [ ] Custom Backoff (TBD)

* See the AWS Blog on [Exponential Backoff and Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/) for more example information on the abovealgorithms

# Installation & Usage

To install backoff, use `go get`:

```
go get github.com/defaltd/backoff
```

Import the `backoff` package into your code using this template:

```go
package main

import (
    "fmt"
    "backoff"
)

func SayHello() {
    service := backoff.New(&backoff.Policy{
        Algorithm: backoff.AlgorithmExponential,
    })

    service.ExecuteAction(func() error { return client.MakeAPICall("Hello World"); })
}
```

# Contributing

Please feel free to contribute by submitting issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.
