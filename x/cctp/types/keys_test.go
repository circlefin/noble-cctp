package types

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestKeys_UsedNonceKey(t *testing.T) {
	tests := []struct {
		name         string
		nonce        uint64
		sourceDomain uint32
		expected     []byte
	}{
		{
			name:         "happy path",
			nonce:        uint64(2),
			sourceDomain: uint32(1),
			expected:     []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, '/'},
		},
		{
			name:         "max value",
			nonce:        math.MaxUint64,
			sourceDomain: math.MaxUint32,
			expected:     []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, '/'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := UsedNonceKey(tt.nonce, tt.sourceDomain)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestKeys_TokenPairKey(t *testing.T) {
	tests := []struct {
		name         string
		remoteDomain uint32
		remoteToken  string
		expected     []byte
	}{
		{
			name:         "happy path",
			remoteDomain: uint32(2),
			remoteToken:  "abc",
			expected:     []byte{0x27, 0xcc, 0x87, 0x7e, 0x1b, 0xc4, 0x72, 0xff, 0xb8, 0x5c, 0xd, 0x32, 0x20, 0x4d, 0xce, 0x7f, 0x5f, 0x7f, 0xfa, 0xfa, 0xf5, 0x91, 0x27, 0x59, 0x44, 0x5a, 0x35, 0xb2, 0xec, 0xc0, 0xb6, 0xda, 0x2f},
		},
		{
			name:         "capitalization doesn't impact",
			remoteDomain: uint32(234),
			remoteToken:  "ABC",
			expected:     TokenPairKey(uint32(234), "abc"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := TokenPairKey(tt.remoteDomain, tt.remoteToken)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestKeys_TokenMessengerKey(t *testing.T) {
	tests := []struct {
		name     string
		domain   uint32
		expected []byte
	}{
		{
			name:     "happy path",
			domain:   uint32(2),
			expected: []byte{0, 0, 0, 2, '/'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := TokenMessengerKey(tt.domain)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
