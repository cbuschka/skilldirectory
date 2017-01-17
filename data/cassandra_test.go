package data

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIDQuery(t *testing.T) {
	options := NewCassandraQueryOptions("key", "id_value", true)
	queryString := " key = id_value"

	if options.Filters[0].query() != queryString {
		t.Errorf("Expecting query to match: %s", queryString)
	}
	fmt.Println(options.Filters[0].query())
}

func TestTextQuery(t *testing.T) {
	options := NewCassandraQueryOptions("key", "text_value", false)
	queryString := " key = 'text_value'"

	if options.Filters[0].query() != queryString {
		t.Errorf("Expecting query to match: %s", queryString)
	}
}

func TestNewOptions(t *testing.T) {
	options1 := CassandraQueryOptions{Filters: []Filter{Filter{"key", "value", true}}}
	options2 := NewCassandraQueryOptions("key", "value", true)

	if !reflect.DeepEqual(options1, options2) {
		t.Error("Expecting NewOptions constructor to match")
	}

	options1.AddFilter("key1", "value1", false)
	options2.Filters = append(options2.Filters, Filter{"key1", "value1", false})

	if !reflect.DeepEqual(options1, options2) {
		t.Error("Expecting AddFilter to match")
	}
}

func TestDeleteSkillChildrean(t *testing.T) {
	table := "links"
	id := ""
	skillID := "1234"
	opts := NewCassandraQueryOptions("skill_ID", skillID, true)
	valid := "DELETE FROM links WHERE skill_ID = 1234;"
	queryString := queryStringHelper(table, id, opts, CassandraConnector{})
	if valid != queryString {
		t.Errorf("Excpecting quesryString to match: %s, %s ", valid, queryString)
	}
}
