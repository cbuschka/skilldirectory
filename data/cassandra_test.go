package data

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Sirupsen/logrus"
)

var logger = logrus.New()

func TestIDQuery(t *testing.T) {
	options := NewQueryOptions("key", "id_value", true)
	queryString := "key = id_value"

	if options.Filters[0].query() != queryString {
		t.Errorf("Expecting query to match: %s", queryString)
	}
	fmt.Println(options.Filters[0].query())
}

func TestTextQuery(t *testing.T) {
	options := NewQueryOptions("key", "text_value", false)
	queryString := "key = 'text_value'"

	if options.Filters[0].query() != queryString {
		t.Errorf("Expecting query to match: %s", queryString)
	}
}

func TestNewOptions(t *testing.T) {
	options1 := QueryOptions{Filters: []Filter{Filter{"key", "value", true}}}
	options2 := NewQueryOptions("key", "value", true)

	if !reflect.DeepEqual(options1, options2) {
		t.Error("Expecting NewOptions constructor to match")
	}

	options1.AddFilter("key1", "value1", false)
	options2.Filters = append(options2.Filters, Filter{"key1", "value1", false})

	if !reflect.DeepEqual(options1, options2) {
		t.Error("Expecting AddFilter to match")
	}
}

func TestDeleteSkillChildren(t *testing.T) {
	table := "links"
	id := ""
	skillID := "1234"
	opts := NewQueryOptions("skill_ID", skillID, true)
	valid := "DELETE FROM links WHERE skill_ID = 1234;"
	logger := logrus.New()
	queryString := makeDeleteQueryStr(table, id, opts, CassandraConnector{Logger: logger})

	if !reflect.DeepEqual(valid, queryString) {
		t.Errorf("Expecting queryString to match: %s, %s", valid, queryString)
	}
}

func TestDeleteNoId(t *testing.T) {
	table := "test_table"
	id := "1234"
	want := "DELETE FROM test_table WHERE id = 1234;"
	got := makeDeleteQueryStr(table, id, QueryOptions{}, CassandraConnector{Logger: logger})

	if got != want {
		t.Errorf("Expected to get: %s ,but got: %s", want, got)
	}
}

func TestDeleteDuplicatedId(t *testing.T) {
	table := "test_table"
	id := "1234"
	opts := NewQueryOptions("id", "1234", true)
	want := "DELETE FROM test_table WHERE id = 1234;"
	got := makeDeleteQueryStr(table, id, opts, CassandraConnector{Logger: logger})

	if got != want {
		t.Errorf("Expected to get: %s ,but got: %s", want, got)
	}
}

func TestDeleteMultipleFilters(t *testing.T) {
	table := "test_table"
	// opts := NewQueryOptions("id", "1234", true)

	opts := NewQueryOptions("hair_color", "red", false)
	opts.AddFilter("name", "Andrew", false)
	want := "DELETE FROM test_table WHERE hair_color = 'red' AND name = 'Andrew';"
	got := makeDeleteQueryStr(table, "", opts, CassandraConnector{Logger: logger})

	// for i := 0; i < len(want); i++ {
	// 	fmt.Print(want[i])
	// 	fmt.Print(" ")
	// 	fmt.Println(got[i])
	// }

	if got != want {
		// fmt.Printf("len(want): %d\nlen(got): %d\n", len(want), len(got))
		// fmt.Println(want)
		// fmt.Println(got)
		t.Errorf("Expected to get: %s ,but got: %s", want, got)
	}
}
