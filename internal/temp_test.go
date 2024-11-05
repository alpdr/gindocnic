package internal

import "testing"

func TestFunc(t *testing.T) {
	if a() != 1 {
		t.Error("expected 1")
	}
}
