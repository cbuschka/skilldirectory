package util

import "fmt"

// FilterMap is a convenience type for applying filters to GORM calls
type FilterMap struct {
	Map map[string]interface{}
}

// NewFilterMap initializes and creates a new FilterMap object with a give key value
func NewFilterMap(key string, value interface{}) *FilterMap {
	f := FilterMap{}
	fMap := make(map[string]interface{})
	fMap[key] = value
	f.Map = fMap
	return &f
}

// Append adds a another key,value pair to a filtermap
func (f *FilterMap) Append(key string, value interface{}) *FilterMap {
	f.Map[key] = value
	return f
}

// WhereQuery converts a key to a 'key IS value' string for GORM where calls
func (f *FilterMap) WhereQuery(key string) string {
	if key == "" {
		return ""
	}
	value := f.Map[key]
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%s IS %v", key, value)
}
