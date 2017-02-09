package util

import "testing"

func TestLogInitNoFlag(t *testing.T) {
	logger := LogInit()
	if logger == nil {
		t.Error("Logger didn't initialize")
	}
}
