package randexc

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Executor struct {
	maxDuration time.Duration
	rand        *rand.Rand
	mu          sync.Mutex
}

type Result struct {
	Error     error
	StartTime time.Time
	EndTime   time.Time
}

type Option func(*Executor) error

func New(maxDuration string, opts ...Option) (*Executor, error) {
	duration, err := time.ParseDuration(maxDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	e := &Executor{
		maxDuration: duration,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Executor) Execute(ctx context.Context, action func() error) error {
	delay := e.randomDelay()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return action()
	}
}

func (e *Executor) ExecuteAsync(ctx context.Context, action func() error) <-chan Result {
	resultChan := make(chan Result, 1)

	go func() {
		startTime := time.Now()
		err := e.Execute(ctx, action)
		resultChan <- Result{
			Error:     err,
			StartTime: startTime,
			EndTime:   time.Now(),
		}
	}()

	return resultChan
}

func (e *Executor) randomDelay() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()
	return time.Duration(e.rand.Int63n(int64(e.maxDuration)))
}

func WithRandSource(source rand.Source) Option {
	return func(e *Executor) error {
		e.rand = rand.New(source)
		return nil
	}
}

func WithMaxDuration(duration string) Option {
	return func(e *Executor) error {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		e.maxDuration = d
		return nil
	}
}
