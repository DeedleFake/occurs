package main

import (
	"bufio"
	"io"
)

type Counter struct {
	Lines map[string]int

	Filters   Filters
	SkipEmpty bool
}

func (c *Counter) Count(r io.Reader) error {
	if c.Lines == nil {
		c.Lines = make(map[string]int)
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := c.Filters.Filter(s.Text())
		if c.SkipEmpty && (line == "") {
			continue
		}

		c.Lines[line]++
	}
	if err := s.Err(); err != nil {
		return err
	}

	return nil
}
