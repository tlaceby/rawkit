package rk

import "testing"

func TestLibrawVersion(t *testing.T) {
	v := Version()
	if v == "" {
		t.Fatal("expected version string, got empty")
	}

	t.Log("libraw version:", v)
}
