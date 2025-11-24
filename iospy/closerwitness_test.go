package iospy

import (
	"errors"
	"testing"

	"github.com/angrifel/unapologetic/ioaux"
)

func TestWitnessCloser(t *testing.T) {
	t.Run("captures successful close", func(t *testing.T) {
		// arrange
		sc := ioaux.CloserFunc(func() error { return nil })

		cw := WitnessCloser(sc).(*closerWitness)

		// act
		closeErr := cw.Close()

		// assert
		if closeErr != nil {
			t.Errorf("expected nil error, got %v", closeErr)
		}

		calls := cw.ObservedCloseCalls()
		if len(calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(calls))
		}

		if calls[0].ResultErr != nil {
			t.Errorf("expected nil error, got %v", calls[0].ResultErr)
		}
		if calls[0].PanicVal != nil {
			t.Errorf("expected nil panic, got %v", calls[0].PanicVal)
		}
	})

	t.Run("captures close with errors", func(t *testing.T) {
		// arrange
		expectedErr := errors.New("read error")
		faultyCloser := ioaux.CloserFunc(func() error { return expectedErr })

		cw := WitnessCloser(faultyCloser).(*closerWitness)

		// act
		closeErr := cw.Close()

		// assert
		if closeErr != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, closeErr)
		}

		calls := cw.ObservedCloseCalls()
		if len(calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(calls))
		}
		if calls[0].ResultErr != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, calls[0].ResultErr)
		}
		if calls[0].PanicVal != nil {
			t.Errorf("expected nil panic, got %v", calls[0].PanicVal)
		}
	})

	t.Run("captures panics and re-panics", func(t *testing.T) {
		// arrange
		expectedPanicVal := "close panic"
		fakeCloser := ioaux.CloserFunc(func() error {
			panic(expectedPanicVal)
		})
		cw := WitnessCloser(fakeCloser).(*closerWitness)

		// act & assert
		var panicVal interface{}
		func() {
			defer func() {
				panicVal = recover()
			}()
			_ = cw.Close()
		}()

		// continue asserting
		if panicVal == nil {
			t.Error("expected panic, but did not panic")
		}
		if panicVal != expectedPanicVal {
			t.Errorf("expected panic value %v, got %v", expectedPanicVal, panicVal)
		}

		calls := cw.ObservedCloseCalls()
		if len(calls) != 1 {
			t.Fatalf("expected 1 call, got %d", len(calls))
		}
		if calls[0].ResultErr != nil {
			t.Errorf("expected nil error, got %v", calls[0].ResultErr)
		}
		if calls[0].PanicVal != panicVal {
			t.Errorf("expected panic value %v, got %v", panicVal, calls[0].PanicVal)
		}
	})
}
