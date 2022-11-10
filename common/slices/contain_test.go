package slices

import "testing"

func TestContainsString(t *testing.T) {
	if ContainsString([]string{"1", "2"}, "3") {
		t.Error(`["1", "2"] > "3"`)
	}

	if !ContainsString([]string{"1", "2"}, "2") {
		t.Error(`["1", "2"] > "2"`)
	}
}
