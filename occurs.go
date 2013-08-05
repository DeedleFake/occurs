package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {
	ts := flag.Bool("ts", false, "Trim whitespace from lines.")
	ic := flag.Bool("ic", false, "Ignore case. (All lines will be converted to lower case.)")
	se := flag.Bool("se", false, "Ignore empty lines.")
	p := flag.Bool("p", false, "Do counting in parallel.")
	flag.Parse()

	filters := make(Filters, 0, 2)
	if *ts {
		filters = append(filters, FilterFunc(strings.TrimSpace))
	}
	if *ic {
		filters = append(filters, FilterFunc(strings.ToLower))
	}

	c := &Counter{
		Filters:   filters,
		SkipEmpty: *se,
	}

	var readers []io.Reader
	if nargs := flag.NArg(); nargs == 0 {
		readers = []io.Reader{os.Stdin}
	} else {
		readers = make([]io.Reader, 0, nargs)
		for _, arg := range flag.Args() {
			if arg == "-" {
				readers = append(readers, os.Stdin)
				continue
			}

			file, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %q because of error: %v\n", arg, err)
				continue
			}
			defer file.Close()

			readers = append(readers, file)
		}
	}

	if *p {
		c.ParallelCount(readers...)
	} else {
		for _, r := range readers {
			err := c.Count(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while counting: %v\n", err)
			}
		}
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', tabwriter.StripEscape)
	for line, num := range c.Lines {
		fmt.Fprintf(tw, "%v\t\xFF%v\xFF\n", num, line)
	}
	tw.Flush()
}
