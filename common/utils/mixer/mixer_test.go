package mixer

import (
	"testing"
)

func TestMixForEach(t *testing.T) {
	MixForEach([][]string{
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
		{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
	}, func(i ...string) error {
		//println(strings.Join(i, ""))
		return nil
	})
}
