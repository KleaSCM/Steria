package metrics

import "testing"

func TestStartProfiling(t *testing.T) {
	prof := StartProfiling()
	if prof == nil {
		t.Fatalf("StartProfiling returned nil")
	}
}

// TODO: Add more unit tests for full coverage
