package main_test

import (
	occurs "."
	"reflect"
	"strings"
	"testing"
)

var testData = []string{
	`this
	is
	a
	test`,

	`so

	Is
	this`,
}

func TestCounter(t *testing.T) {
	c := &occurs.Counter{
		Filters: occurs.Filters{
			occurs.FilterFunc(strings.TrimSpace),
			occurs.FilterFunc(strings.ToLower),
		},

		SkipEmpty: true,
	}

	for i, data := range testData {
		err := c.Count(strings.NewReader(data))
		if err != nil {
			t.Fatalf("Error when counting %v: %v", i, err)
		}
	}

	expected := map[string]int{
		"this": 2,
		"is":   2,
		"a":    1,
		"test": 1,
		"so":   1,
	}

	t.Logf("%q", c.Lines)
	t.Logf("%q", expected)

	if !reflect.DeepEqual(c.Lines, expected) {
		t.Fatalf("%q != %q", c.Lines, expected)
	}
}
