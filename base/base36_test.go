package base

import (
	"testing"
)

var tests = []uint64{123456, 233333, 123456789, 2821109907455, 2821109907455}

var encoded = []string{"2n9c", "501h", "21i3v9", "zzzzzzzz", "ZZZZZZZZ"}

func TestEncode(t *testing.T) {
	for i, tt := range tests {
		t.Run("Encode uint64 to base36 string", func(t *testing.T) {
			if got := Encode(tt); got != encoded[i] {
				t.Errorf("Encode() = %v, want %v", got, encoded[i])
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for i, tt := range encoded {
		t.Run("Decode base36 string to uint64", func(t *testing.T) {
			if got := Decode(tt); got != tests[i] {
				t.Errorf("Decode() = %v, want %v", got, tests[i])
			}
		})
	}
}
