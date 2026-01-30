package core

import "testing"

func TestLibrawVersion(t *testing.T) {
	v := LibrawVersion()
	if v == "" {
		t.Fatal("expected version string, got empty")
	}

	t.Log("libraw version:", v)
}
