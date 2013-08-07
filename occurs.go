package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

var c *Counter

func countSeq() {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

	for _, arg := range args {
		if arg == "-" {
			err := c.Count(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while counting from stdin: %v\n", err)
			}
			continue
		}

		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %q because of error: %v\n", arg, err)
			continue
		}

		err = c.Count(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error occured while counting from %q: %v\n", arg, err)
		}

		file.Close()
	}
}

func countParallel() {
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

			readers = append(readers, file)
		}
	}

	c.ParallelCount(readers...)
}

var (
	ts = flag.Bool("ts", false, "Trim whitespace from lines.")
	ic = flag.Bool("ic", false, "Ignore case. (All lines will be converted to lower case.)")
	se = flag.Bool("se", false, "Ignore empty lines.")

	p = flag.Bool("p", false, "Do counting in parallel. (This causes all files to be opened at once.)")

	cols = flag.Bool("cols", false, "Output in nicely aligned columns.")
)

func main() {
	flag.Parse()

	filters := make(Filters, 0, 2)
	if *ts {
		filters = append(filters, FilterFunc(strings.TrimSpace))
	}
	if *ic {
		filters = append(filters, FilterFunc(strings.ToLower))
	}

	c = &Counter{
		Filters:   filters,
		SkipEmpty: *se,
	}

	if *p {
		countParallel()
	} else {
		countSeq()
	}

	if *cols {
		tw := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', tabwriter.StripEscape)
		for line, num := range c.Lines {
			fmt.Fprintf(tw, "%v\t\xFF%v\xFF\n", num, line)
		}
		tw.Flush()
	} else {
		for line, num := range c.Lines {
			fmt.Printf("%v %v\n", num, line)
		}
	}
}
