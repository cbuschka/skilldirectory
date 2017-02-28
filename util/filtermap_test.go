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
