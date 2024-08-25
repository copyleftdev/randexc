# randexc: Random Execution Library Design Document

## 1. Introduction

### 1.1 Purpose
The `randexc` library is designed to provide a simple and flexible way to execute actions at random times within a specified duration. This library aims to solve the problem of scheduling tasks with a degree of randomness, which can be useful in various scenarios such as load testing, simulating real-world events, or implementing exponential backoff strategies.

### 1.2 Scope
This library will provide a Go package that allows users to:
- Define a maximum time frame for execution
- Execute actions randomly within that time frame
- Support context-based cancellation and timeouts
- Provide both synchronous and asynchronous execution options

## 2. System Architecture

### 2.1 High-Level Architecture
The `randexc` library will be implemented as a single Go package with a main `Executor` type that handles the scheduling and execution of random timed actions.

### 2.2 Components
1. Executor: The main struct that manages the random execution logic
2. Options: A set of functional options to configure the Executor
3. Result: A struct to represent the outcome of an execution

## 3. Detailed Design

### 3.1 Executor

```go
type Executor struct {
    maxDuration time.Duration
    rand        *rand.Rand // For testability
}

func New(maxDuration string, opts ...Option) (*Executor, error)
func (e *Executor) Execute(ctx context.Context, action func() error) error
func (e *Executor) ExecuteAsync(ctx context.Context, action func() error) <-chan Result
```

#### 3.1.1 New
Creates a new Executor with the specified maximum duration and applies any provided options.

#### 3.1.2 Execute
Executes the given action synchronously within the random time frame.

#### 3.1.3 ExecuteAsync
Executes the given action asynchronously within the random time frame and returns a channel for the result.

### 3.2 Options

```go
type Option func(*Executor) error

func WithRandSource(source rand.Source) Option
func WithMaxAttempts(attempts int) Option
```

#### 3.2.1 WithRandSource
Allows setting a custom random source for testability.

#### 3.2.2 WithMaxAttempts
Sets the maximum number of execution attempts in case of failures.

### 3.3 Result

```go
type Result struct {
    Error     error
    StartTime time.Time
    EndTime   time.Time
}
```

Represents the outcome of an execution, including any error, start time, and end time.

## 4. Interfaces

### 4.1 Main Interface

```go
type Randomizer interface {
    Execute(ctx context.Context, action func() error) error
    ExecuteAsync(ctx context.Context, action func() error) <-chan Result
}
```

This interface defines the main methods that the Executor should implement, allowing for easy mocking in tests.

## 5. Error Handling

The library will use custom error types to provide more context about failures:

```go
var (
    ErrInvalidDuration = errors.New("invalid duration")
    ErrMaxAttemptsReached = errors.New("maximum execution attempts reached")
    ErrContextCancelled = errors.New("execution cancelled due to context cancellation")
)
```

## 6. Concurrency and Thread Safety

The `Executor` will be designed to be safe for concurrent use. The internal `rand.Rand` will be protected by a mutex if necessary.

## 7. Performance Considerations

- The library will use efficient random number generation to determine execution times.
- For very short durations, we'll ensure that the overhead of random number generation doesn't significantly impact the timing.

## 8. Testing Strategy

- Unit tests will cover all public methods and error cases.
- We'll use dependency injection (WithRandSource option) to make random behaviors deterministic in tests.
- Benchmarks will be created to ensure performance for various duration ranges.

## 9. Documentation

- Godoc-compatible comments will be used for all exported types and functions.
- A usage guide with examples will be provided in the `docs` directory.
- The README will include quick start examples and installation instructions.

## 10. Future Enhancements

- Support for periodic random executions
- Custom distribution functions for randomness (e.g., normal distribution)
- Integration with popular scheduling libraries

## 11. Conclusion

The `randexc` library aims to provide a robust and flexible solution for random-timed executions in Go. Its simple API, combined with powerful features like context support and async execution, will make it a valuable tool for developers needing controlled randomness in their applications.