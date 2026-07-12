# ⏳ lesiw.io/wait 

[![Go Reference](https://pkg.go.dev/badge/lesiw.io/wait.svg)](https://pkg.go.dev/lesiw.io/wait)
[![CI](https://github.com/lesiw/wait/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/lesiw/wait/actions/workflows/main.yml)
[![Release](https://img.shields.io/github/v/tag/lesiw/wait?sort=semver&label=release)](https://github.com/lesiw/wait/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lesiw/wait)](go.mod)
[![Discord](https://img.shields.io/discord/1145827224516300971?logo=discord&logoColor=white&color=5865F2&label=discord)](https://lesiw.dev/discord)
[![License](https://img.shields.io/github/license/lesiw/wait)](../LICENSE)

Package wait provides waiting strategies. A Waiter returns a channel
that fires when the caller may proceed.

## Features

* **Channel-native.** `Wait` returns a channel; compose with
  `select` and `context.Context`.
* **Two built-in strategies.** `Decay` for supervisors, `Jitter` for retries.
* **Zero dependencies.** Standard library only.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "lesiw.io/wait"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    var (
        start = time.Now()
        since = start
        tries = 0
    )
    fn := func() error {
        tries++
        now := time.Now()
        fmt.Printf("attempt=%d elapsed=%v delay=%v\n",
            tries,
            now.Sub(start).Round(time.Millisecond),
            now.Sub(since).Round(time.Millisecond),
        )
        since = now
        if tries < 5 {
            return fmt.Errorf("bang!")
        }
        return nil
    }
    wtr := wait.Jitter(2 * time.Second)
    for fn() != nil {
        select {
        case <-wtr.Wait():
        case <-ctx.Done():
            return
        }
    }
    fmt.Println("succeeded after", tries, "tries")
}
```

[▶️ Run this example on the Go Playground](https://go.dev/play/p/tLr7fADNMm1)

## Strategies

| Waiter          | Behavior                                                                                            |
|-----------------|-----------------------------------------------------------------------------------------------------|
| `Decay(limit)`  | Delay increases under rapid calls and decays with idle time. Good for supervisor restart loops.     |
| `Jitter(limit)` | Random `[0, upper)`; `upper` doubles from `limit/1024` to `limit`. Good for retry pacing.           |

## Motivation

Most waiting strategies are simple to understand but subtle to implement.
They're also usually a part of the error path,
which tends to suffer in terms of test-coverage and real-world observation.

So for that reason I think it makes sense to group them together
as implementations to a simple Waiter() interface.

Tunables are deliberately kept to a minimum to avoid having to
double-check math at every callsite.
Typically the only interesting question is
"What's the maximum delay you will tolerate?"
and so the Waiters in this package are built around that question.

## Further Reading

- [Suture v4 Spec](https://pkg.go.dev/github.com/thejerf/suture/v4#Spec)
  — the failure-decay counter `Decay` is inspired by.
- [Exponential Backoff And Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
  — Marc Brooker's 2015 AWS engineering post that defines and
  analyzes the Full Jitter algorithm `Jitter` uses.
