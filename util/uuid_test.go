package util

import (
	"testing"
	"time"
)

func TestNewId(t *testing.T) {
	id1 := NewID()
	time.Sleep(1)
	id2 := NewID()
	if id1 == id2 {
		t.Error("Random id failed")
	}

}

func TestIdLength(t *testing.T) {
	id := NewID()
	if len(id) != 36 {
		t.Errorf("ID length incorrect")
	}
}
