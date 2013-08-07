package main

import (
	"bufio"
	"io"
	"sync"
)

type Counter struct {
	Lines map[string]int

	Filters   Filters
	SkipEmpty bool
}

func (c *Counter) init() {
	if c.Lines == nil {
		c.Lines = make(map[string]int)
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

func (c *Counter) Count(r io.Reader) error {
	c.init()

	return c.lowerCount(r, func(line string) {
		c.Lines[line]++
	})
}

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
