package randexc

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		maxDuration string
		wantErr     bool
	}{
		{"Valid duration", "1h", false},
		{"Invalid duration", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.maxDuration)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutor_Execute(t *testing.T) {
	executor, _ := New("1s", WithRandSource(rand.NewSource(0)))

	t.Run("Successful execution", func(t *testing.T) {
		err := executor.Execute(context.Background(), func() error {
			return nil
		})
		if err != nil {
			t.Errorf("Execute() error = %v, want nil", err)
		}
	})

	t.Run("Action error", func(t *testing.T) {
		expectedErr := errors.New("action error")
		err := executor.Execute(context.Background(), func() error {
			return expectedErr
		})
		if err != expectedErr {
			t.Errorf("Execute() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("Context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := executor.Execute(ctx, func() error {
			return nil
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Execute() error = %v, want context.Canceled", err)
		}
	})
}

func TestExecutor_ExecuteAsync(t *testing.T) {
	executor, _ := New("1s", WithRandSource(rand.NewSource(0)))

	t.Run("Successful async execution", func(t *testing.T) {
		resultChan := executor.ExecuteAsync(context.Background(), func() error {
			return nil
		})
		result := <-resultChan
		if result.Error != nil {
			t.Errorf("ExecuteAsync() error = %v, want nil", result.Error)
		}
		if result.StartTime.IsZero() || result.EndTime.IsZero() {
			t.Error("ExecuteAsync() start or end time is zero")
		}
	})

	t.Run("Async action error", func(t *testing.T) {
		expectedErr := errors.New("async action error")
		resultChan := executor.ExecuteAsync(context.Background(), func() error {
			return expectedErr
		})
		result := <-resultChan
		if result.Error != expectedErr {
			t.Errorf("ExecuteAsync() error = %v, want %v", result.Error, expectedErr)
		}
	})
}

func TestWithRandSource(t *testing.T) {
	source := rand.NewSource(42)
	executor, _ := New("1s", WithRandSource(source))

	delay1 := executor.randomDelay()
	delay2 := executor.randomDelay()

	if delay1 == delay2 {
		t.Error("WithRandSource() did not set custom source correctly")
	}
}

func TestWithMaxDuration(t *testing.T) {
	executor, _ := New("1s", WithMaxDuration("2s"))

	if executor.maxDuration != 2*time.Second {
		t.Errorf("WithMaxDuration() maxDuration = %v, want 2s", executor.maxDuration)
	}
}

func BenchmarkExecutor_Execute(b *testing.B) {
	executor, _ := New("1s")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = executor.Execute(ctx, func() error {
			return nil
		})
	}
}
