// Package testaux provides auxiliary functions for testing.
package testaux

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"time"
)

// ErrTimeoutWaitingForConnect is returned when WaitForConnectAccept times out.
var ErrTimeoutWaitingForConnect = errors.New("timed out waiting for TCP connection")

// SliceEqual is a helper function for comparing slices of comparable values.
func SliceEqual[T comparable](a, b []T) bool { return slices.Equal(a, b) }

// HeaderEqual compares two http.Headers for equality.
func HeaderEqual(a, b http.Header) bool {
	if len(a) != len(b) {
		return false
	}

	for k, aValue := range a {
		bValue, keyExistsInB := b[k]
		if !keyExistsInB {
			return false
		}

		if !slices.Equal(aValue, bValue) {
			return false
		}
	}

	return true
}

// WaitForConnectAccept waits until a TCP connection can be established to the given address.
func WaitForConnectAccept(address string) error {
	const (
		waitTimeout = 5 * time.Millisecond
		maxWaitTime = 1 * time.Second
	)

	startTime := time.Now()
	dialer := net.Dialer{Timeout: waitTimeout} //nolint:exhaustruct // Partial initialization is fine in this context

	for {
		conn, connErr := dialer.DialContext(context.Background(), "tcp", address)
		if connErr == nil {
			return conn.Close()
		} else if time.Since(startTime) > maxWaitTime {
			return fmt.Errorf("%w: %s", ErrTimeoutWaitingForConnect, address)
		}
	}
}
