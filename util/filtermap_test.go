package util

import (
	"reflect"
	"testing"
)

func TestFilterMap(t *testing.T) {
	fmap := &FilterMap{}
	fmapMap := make(map[string]interface{})
	fmap.Map = fmapMap
	fmap.Map["a"] = "b"
	fmap2 := NewFilterMap("a", "b")
	if !reflect.DeepEqual(fmap, fmap2) {
		t.Errorf("Map: %v doesn't match Map2: %v", fmap, fmap2)
	}
}

func TestAppend(t *testing.T) {
	filterMap := NewFilterMap("a", "b").Append("c", "d")
	if filterMap.Map["c"] != "d" {
		t.Errorf("Map[c] is: %s should be 'd'", filterMap.Map["c"])
	}
}

func TestWhereNoKey(t *testing.T) {
	fm := NewFilterMap("key", "value")
	if response := fm.WhereQuery(""); response != "" {
		t.Errorf("Expected '' got: %s", response)
	}
}

func TestWhereMissingKey(t *testing.T) {
	fm := NewFilterMap("key", "value")
	if response := fm.WhereQuery("key2"); response != "" {
		t.Errorf("Expected '' got: %s", response)
	}
}

func TestWhere(t *testing.T) {
	fm := NewFilterMap("key", "value")
	if response := fm.WhereQuery("key"); response != "key IS value" {
		t.Errorf("Expected '' got: %s", response)
	}
}
