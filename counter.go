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

func (c *Counter) Count(r io.Reader) error {
	c.init()

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

			s := bufio.NewScanner(r)
			for s.Scan() {
				line := c.Filters.Filter(s.Text())
				if c.SkipEmpty && (line == "") {
					continue
				}

				lineC <- line
			}
		}(r)
	}
	wg.Wait()

	close(lineC)
	<-done
}
