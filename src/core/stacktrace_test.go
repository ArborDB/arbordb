package core

import (
	"strings"
	"testing"
)

func TestStacktrace(t *testing.T) {
	stacktrace := NewStacktrace()
	msg := stacktrace.Error()
	if !strings.Contains(msg, "TestStacktrace") {
		t.Fatalf("got %s", msg)
	}
}
