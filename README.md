# randexc: Random Execution Library for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/copyleftdev/randexc)](https://goreportcard.com/report/github.com/copyleftdev/randexc)
[![GoDoc](https://godoc.org/github.com/copyleftdev/randexc?status.svg)](https://godoc.org/github.com/copyleftdev/randexc)
[![Build Status](https://github.com/copyleftdev/randexc/workflows/Go/badge.svg)](https://github.com/copyleftdev/randexc/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`randexc` is a powerful and flexible Go library for executing actions at random times within a specified duration. It provides both synchronous and asynchronous execution options, making it suitable for a wide range of applications including load testing, simulating real-world events, and implementing sophisticated retry mechanisms.

## Features

- 🕒 Execute actions within a random time frame
- 🔄 Support for both synchronous and asynchronous execution
- 🛠 Configurable through functional options
- 🧪 Easy to test with custom random sources
- 📦 Lightweight with no external dependencies

## Installation

```bash
go get github.com/copyleftdev/randexc
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/copyleftdev/randexc/pkg/randexc"
)

func main() {
    executor, err := randexc.New("1m")
    if err != nil {
        log.Fatalf("Failed to create executor: %v", err)
    }

    err = executor.Execute(context.Background(), func() error {
        fmt.Println("Action executed at a random time within 1 minute!")
        return nil
    })

    if err != nil {
        log.Printf("Execution failed: %v", err)
    }
}
```

## Use Cases

### Load Testing

Simulate realistic user behavior in load tests by executing actions at random intervals:

```go
func performRequest() error {
    // Simulate HTTP request
    time.Sleep(100 * time.Millisecond)
    return nil
}

executor, _ := randexc.New("5s")

for i := 0; i < 100; i++ {
    go func() {
        err := executor.Execute(context.Background(), performRequest)
        if err != nil {
            log.Printf("Request failed: %v", err)
        }
    }()
}
```

### Jittered Exponential Backoff

Implement a jittered exponential backoff strategy for retrying operations:

```go
func retryWithBackoff(operation func() error) error {
    baseDelay := 100 * time.Millisecond
    maxDelay := 10 * time.Second
    maxAttempts := 5

    for attempt := 0; attempt < maxAttempts; attempt++ {
        jitteredDelay := time.Duration(float64(baseDelay) * (1 + rand.Float64()))
        if jitteredDelay > maxDelay {
            jitteredDelay = maxDelay
        }

        executor, _ := randexc.New(jitteredDelay.String())
        err := executor.Execute(context.Background(), operation)
        if err == nil {
            return nil
        }

        baseDelay *= 2
    }

    return fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
```

### Simulating Real-world Events

Create more realistic simulations by introducing randomness in event timing:

```go
func simulateIoTDevice(deviceID string) {
    executor, _ := randexc.New("1h")

    for {
        executor.Execute(context.Background(), func() error {
            temperature := 20 + rand.Float64()*10
            humidity := 30 + rand.Float64()*20
            fmt.Printf("Device %s: Temp: %.2f°C, Humidity: %.2f%%\n", deviceID, temperature, humidity)
            return nil
        })
    }
}

// Simulate multiple IoT devices
for i := 0; i < 5; i++ {
    go simulateIoTDevice(fmt.Sprintf("device_%d", i))
}
```

## Documentation

For more detailed documentation and advanced usage examples, please see the [usage guide](docs/usage.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
