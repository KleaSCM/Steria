package utils

import "testing"

func TestShouldIgnore(t *testing.T) {
	patterns := []IgnorePattern{{Pattern: "*.log", IsDir: false}, {Pattern: "temp", IsDir: true}}
	if !ShouldIgnore("foo.log", patterns) {
		t.Errorf("ShouldIgnore failed for foo.log")
	}
	if ShouldIgnore("main.go", patterns) {
		t.Errorf("ShouldIgnore incorrectly ignored main.go")
	}
}

// TODO: Add more unit tests for full coverage
