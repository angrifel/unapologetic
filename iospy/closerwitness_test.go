package iospy

import (
	"errors"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/ioaux"
	"github.com/angrifel/unapologetic/internal/require"
)

func TestWitnessCloser(t *testing.T) {
	t.Run("captures successful close", func(t *testing.T) {
		// arrange
		sc := ioaux.CloserFunc(func() error { return nil })

		cw := WitnessCloser(sc).(*closerWitness)

		// act
		closeErr := cw.Close()

		// assert
		assert.IsNil(t, closeErr)

		calls := cw.ObservedCloseCalls()
		require.Equal(t, 1, len(calls))

		assert.IsNil(t, calls[0].ResultErr)
		assert.IsNil(t, calls[0].PanicVal)
	})

	t.Run("captures close with errors", func(t *testing.T) {
		// arrange
		expectedErr := errors.New("read error")
		faultyCloser := ioaux.CloserFunc(func() error { return expectedErr })

		cw := WitnessCloser(faultyCloser).(*closerWitness)

		// act
		closeErr := cw.Close()

		// assert
		assert.Equal(t, expectedErr, closeErr)

		calls := cw.ObservedCloseCalls()
		require.Equal(t, 1, len(calls))
		assert.Equal(t, expectedErr, calls[0].ResultErr)
		assert.IsNil(t, calls[0].PanicVal)
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
			t.Fatal("expected panic, but did not panic")
		}
		assert.Equal[any](t, expectedPanicVal, panicVal)

		calls := cw.ObservedCloseCalls()
		require.Equal(t, 1, len(calls))
		assert.IsNil(t, calls[0].ResultErr)
		assert.Equal(t, panicVal, calls[0].PanicVal)
	})
}
