package testaux

import (
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/angrifel/unapologetic/internal/assert"
	"github.com/angrifel/unapologetic/internal/require"
)

func TestSliceEqual(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want bool
	}{
		{
			name: "equal slices",
			a:    []int{1, 2, 3},
			b:    []int{1, 2, 3},
			want: true,
		},
		{
			name: "different lengths",
			a:    []int{1, 2},
			b:    []int{1, 2, 3},
			want: false,
		},
		{
			name: "different values",
			a:    []int{1, 2, 3},
			b:    []int{1, 2, 4},
			want: false,
		},
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "both empty",
			a:    []int{},
			b:    []int{},
			want: true,
		},
		{
			name: "one nil one empty",
			a:    nil,
			b:    []int{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SliceEqual(tt.a, tt.b))
		})
	}
}

func TestSliceEqualStrings(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "equal string slices",
			a:    []string{"foo", "bar"},
			b:    []string{"foo", "bar"},
			want: true,
		},
		{
			name: "different string slices",
			a:    []string{"foo", "bar"},
			b:    []string{"foo", "baz"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SliceEqual(tt.a, tt.b))
		})
	}
}

func TestHeaderEqual(t *testing.T) {
	tests := []struct {
		name string
		a    http.Header
		b    http.Header
		want bool
	}{
		{
			name: "equal headers",
			a: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"text/html"},
			},
			b: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"text/html"},
			},
			want: true,
		},
		{
			name: "different lengths",
			a: http.Header{
				"Content-Type": []string{"application/json"},
			},
			b: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"text/html"},
			},
			want: false,
		},
		{
			name: "missing key in b",
			a: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"text/html"},
			},
			b: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer secret-token"},
			},
			want: false,
		},
		{
			name: "different values",
			a: http.Header{
				"Content-Type": []string{"application/json"},
			},
			b: http.Header{
				"Content-Type": []string{"text/html"},
			},
			want: false,
		},
		{
			name: "different value counts",
			a: http.Header{
				"Accept": []string{"text/html", "application/json"},
			},
			b: http.Header{
				"Accept": []string{"text/html"},
			},
			want: false,
		},
		{
			name: "multiple values in different order",
			a: http.Header{
				"Accept": []string{"text/html", "application/json"},
			},
			b: http.Header{
				"Accept": []string{"application/json", "text/html"},
			},
			want: false,
		},
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "both empty",
			a:    http.Header{},
			b:    http.Header{},
			want: true,
		},
		{
			name: "one nil one empty",
			a:    nil,
			b:    http.Header{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, HeaderEqual(tt.a, tt.b))
		})
	}
}

func TestWaitForConnectAccept(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		// Start a listener on a random available port
		listener, listerErr := net.Listen("tcp", "127.0.0.1:0")
		require.IsNil(t, listerErr)
		t.Cleanup(func() { _ = listener.Close() })

		// Get the actual address that was bound
		address := listener.Addr().String()

		// Accept connections in the background
		go func() {
			for {
				conn, connErr := listener.Accept()
				if connErr != nil {
					return
				}
				_ = conn.Close()
			}
		}()

		// Test that WaitForConnectAccept succeeds
		assert.IsNil(t, WaitForConnectAccept(address))
	})

	t.Run("timeout on unavailable port", func(t *testing.T) {
		// Use an address that won't accept connections
		// Find an available port and then don't listen on it
		listener, listenerErr := net.Listen("tcp", "127.0.0.1:0")
		require.IsNil(t, listenerErr)
		address := listener.Addr().String()
		_ = listener.Close() // Close immediately so nothing is listening

		// Test that WaitForConnectAccept times out
		err := WaitForConnectAccept(address)
		assert.IsNotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrTimeoutWaitingForConnect))
	})

	t.Run("invalid address", func(t *testing.T) {
		// Test with an invalid address format
		err := WaitForConnectAccept("invalid:address:format")
		assert.IsNotNil(t, err)

		// Should timeout since it can't connect
		assert.Equal(t, true, errors.Is(err, ErrTimeoutWaitingForConnect))
	})
}
