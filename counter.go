package main

import (
	"bufio"
	"io"
	"sync"
)

// A Counter is used for counting lines. To reset the counter, simply
// set Lines to nil. The zero value for a counter is usable.
type Counter struct {
	// Lines is a map of the lines. The key for the map is the line, and
	// the value is the number of times that that line has occurred.
	Lines map[string]uint

	// Filters is a list of Filterers to use to filter lines.
	Filters Filters

	// If SkipEmpty is true empty lines are not counted. This is checked
	// after the Filters are run, so if this is true and the Filters
	// output an empty line, the line is ignored.
	SkipEmpty bool
}

func (c *Counter) init() {
	if c.Lines == nil {
		c.Lines = make(map[string]uint)
	}
}

func (c *Counter) lowerCount(r io.Reader, f func(string)) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := c.Filters.Filter(s.Text())
		if c.SkipEmpty && (line == "") {
			continue
		}

		f(line)
	}

	return s.Err()
}

// Count reads from r until it encounters an error, counting all the
// lines it reads. It returns the error encountered, unless the error
// is io.EOF.
func (c *Counter) Count(r io.Reader) error {
	c.init()

	return c.lowerCount(r, func(line string) {
		c.Lines[line]++
	})
}

// ParallelCount counts from all the io.Readers in r in parallel.
// Because of the nature of this method it does not return an error,
// even if one is encountered. This method is not asynchronous, and
// does not return until it has finished counting from all the
// io.Readers.
func (c *Counter) ParallelCount(r ...io.Reader) {
	c.init()

	lineC := make(chan string, 1024)
	done := make(chan bool)

	go func() {
		for line := range lineC {
			c.Lines[line]++
		}

		done <- true
	}()

	var wg sync.WaitGroup
	for _, r := range r {
		wg.Add(1)
		go func(r io.Reader) {
			defer wg.Done()

			c.lowerCount(r, func(line string) {
				lineC <- line
			})
		}(r)
	}
	wg.Wait()

	close(lineC)
	<-done
}
