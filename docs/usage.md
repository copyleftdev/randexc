# randexc Usage Guide

`randexc` is a Go library for executing actions at random times within a specified duration. This guide will help you get started with using `randexc` in your projects.

## Table of Contents

1. [Installation](#installation)
2. [Basic Usage](#basic-usage)
3. [Advanced Features](#advanced-features)
   - [Asynchronous Execution](#asynchronous-execution)
   - [Custom Random Source](#custom-random-source)
   - [Changing Maximum Duration](#changing-maximum-duration)
4. [Best Practices](#best-practices)
5. [Examples](#examples)
6. [Troubleshooting](#troubleshooting)

## Installation

To install `randexc`, use the following command:

```bash
go get github.com/yourusername/randexc
```

Then import it in your Go code:

```go
import "github.com/yourusername/randexc/pkg/randexc"
```

## Basic Usage

Here's a simple example of how to use `randexc`:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yourusername/randexc/pkg/randexc"
)

func main() {
    // Create a new Executor with a maximum duration of 1 minute
    executor, err := randexc.New("1m")
    if err != nil {
        log.Fatalf("Failed to create executor: %v", err)
    }

    // Execute an action
    err = executor.Execute(context.Background(), func() error {
        fmt.Println("Action executed!")
        return nil
    })

    if err != nil {
        log.Printf("Execution failed: %v", err)
    }
}
```

This will execute the action (printing "Action executed!") at a random time within the next minute.

## Advanced Features

### Asynchronous Execution

For non-blocking execution, use `ExecuteAsync`:

```go
resultChan := executor.ExecuteAsync(context.Background(), func() error {
    fmt.Println("Async action executed!")
    return nil
})

// Do other work...

result := <-resultChan
if result.Error != nil {
    log.Printf("Async execution failed: %v", result.Error)
} else {
    fmt.Printf("Async execution completed. Start: %v, End: %v\n", result.StartTime, result.EndTime)
}
```

### Custom Random Source

You can provide a custom random source for deterministic behavior (useful in tests):

```go
import "math/rand"

executor, err := randexc.New("1m", randexc.WithRandSource(rand.NewSource(42)))
```

### Changing Maximum Duration

You can change the maximum duration after creation:

```go
executor, _ := randexc.New("1m", randexc.WithMaxDuration("2h"))
```

## Best Practices

1. Always check for errors when creating a new `Executor` and when calling `Execute` or `ExecuteAsync`.
2. Use appropriate durations for your use case. Very short durations might not provide enough randomness, while very long durations might delay your action too much.
3. For long-running or potentially expensive operations, use context cancellation to allow for graceful shutdowns.
4. In tests, use a custom random source for reproducibility.

## Examples

### Load Testing

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

### Exponential Backoff with Randomness

```go
func exponentialBackoff(baseDelay time.Duration, maxAttempts int) {
    for attempt := 0; attempt < maxAttempts; attempt++ {
        executor, _ := randexc.New(baseDelay.String())
        err := executor.Execute(context.Background(), performAction)
        if err == nil {
            return
        }
        baseDelay *= 2
    }
}
```

## Troubleshooting

- If actions are not executing, ensure the provided duration is not too long.
- For issues with randomness, check if a custom random source is being used unintentionally.
- If experiencing unexpected delays, verify that the context is not being cancelled prematurely.

For more help, please open an issue on the GitHub repository.