package util

import (
	"fmt"
	"testing"
)

func Test_Utob(t *testing.T) {
	tests := []struct {
		name string
		args uint64
		want [8]byte
	}{
		// don't need test
		{"test 0 ", uint64(0), [8]byte(uint8(0))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Utob(tt.args)
			fmt.Println(got)
		})
	}
}
