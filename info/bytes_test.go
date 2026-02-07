package info

import (
	"math"
	"testing"
)

func TestBytes(t *testing.T) {
	btests := []struct {
		value  Bytes
		output string
	}{
		{value: Bytes(512), output: "512 B"},
		{value: MBytes(1), output: "1.00 MB"},
		{value: GBytes(1), output: "1.00 GB"},
		{value: Bytes(2.5 * 1024 * 1024), output: "2.50 MB"},
		{value: Bytes(math.Round(3.7 * 1024 * 1024 * 1024)), output: "3.70 GB"},
	}

	for _, bt := range btests {
		t.Run(bt.output, func(t *testing.T) {
			got := bt.value.String()
			if got != bt.output {
				t.Errorf("%#v got %s want %s", bt.value, got, bt.output)
			}
		})
	}
}
