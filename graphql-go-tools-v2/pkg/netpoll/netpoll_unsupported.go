//go:build windows || wasm
// +build windows wasm

package netpoll

import (
	"time"
)

// NewPoller creates a new poll based connection implementation.
func NewPoller(connBufferSize int, _ time.Duration) (Poller, error) {
	return nil, ErrUnsupported
}
