package main_test

import (
	occurs "."
	"io"
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

	expected := map[string]int{
		"this": 2,
		"is":   2,
		"a":    1,
		"test": 1,
		"so":   1,
	}
	t.Logf("Expecting: %q", expected)

	for i, data := range testData {
		err := c.Count(strings.NewReader(data))
		if err != nil {
			t.Fatalf("Error when counting %v: %v", i, err)
		}
	}

	t.Logf("Sequential: %q", c.Lines)

	if !reflect.DeepEqual(c.Lines, expected) {
		t.Fatalf("Sequential: %q != %q", c.Lines, expected)
	}

	c.Lines = nil

	readers := make([]io.Reader, 0, len(testData))
	for _, data := range testData {
		readers = append(readers, strings.NewReader(data))
	}
	c.ParallelCount(readers...)

	t.Logf("Parallel: %q", c.Lines)

	if !reflect.DeepEqual(c.Lines, expected) {
		t.Fatalf("Parallel: %q != %q", c.Lines, expected)
	}
}
