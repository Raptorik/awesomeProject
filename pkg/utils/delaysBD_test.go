package utils

import (
	"errors"
	"testing"
	"time"
)

func TestDoWithTries(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		attempts := 3
		delay := 100 * time.Millisecond
		fn := func() error {
			return nil
		}

		err := DoWithTries(fn, attempts, delay)
		if err != nil {
			t.Fatalf("DoWithTries() should not return an error when the function is successful: %v", err)
		}
	})

	t.Run("failed execution", func(t *testing.T) {
		attempts := 3
		delay := 100 * time.Millisecond
		fn := func() error {
			return errors.New("test error")
		}

		err := DoWithTries(fn, attempts, delay)
		if err == nil {
			t.Fatal("DoWithTries() should return an error when the function is never successful")
		}
	})

	t.Run("successful execution after retries", func(t *testing.T) {
		attempts := 3
		delay := 100 * time.Millisecond
		counter := 0
		fn := func() error {
			counter++
			if counter == attempts {
				return nil
			}
			return errors.New("test error")
		}

		err := DoWithTries(fn, attempts, delay)
		if err != nil {
			t.Fatalf("DoWithTries() should not return an error when the function is successful after retries: %v", err)
		}
	})
}
