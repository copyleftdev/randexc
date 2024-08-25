// Package randexc provides functionality for executing actions at random times within a specified duration.
//
// The main type in this package is Executor, which allows scheduling and performing
// actions within a given timeframe. It supports both synchronous and asynchronous
// execution modes, making it suitable for various use cases where randomized timing is desired.
//
// Key features:
//   - Random execution within a specified maximum duration
//   - Support for context-based cancellation
//   - Synchronous and asynchronous execution options
//   - Configurable through functional options
//
// Basic usage:
//
//	executor, err := randexc.New("1h")
//	if err != nil {
//		// handle error
//	}
//
//	err = executor.Execute(context.Background(), func() error {
//		// perform action
//		return nil
//	})
//
// For asynchronous execution:
//
//	resultChan := executor.ExecuteAsync(context.Background(), func() error {
//		// perform action
//		return nil
//	})
//
//	result := <-resultChan
//	if result.Error != nil {
//		// handle error
//	}
//
// This package is useful for scenarios such as load testing, simulating real-world events,
// or implementing exponential backoff strategies with a random component.
package randexc
