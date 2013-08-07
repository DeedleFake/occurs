// occurs is a simple program which counts the number of times unique
// lines occur in a given set of files. For example, given the files
//
//	example1.txt
//	this
//	is
//	an
//	example
//
// and
//
//	example2.txt
//	this
//	is
//	also
//	an
//	example
//
// running
//
//	occurs example1.txt example2.txt
//
// will output
//
//	2 this
//	2 is
//	2 an
//	1 also
//	2 example
//
// Note that the order of the output lines is unspecified, and may
// output in a different order when run on the same files multiple
// times. To enforce an order to the output, you can pipe the output
// into another program. For example, the Unix sort command; running
//
//	occurs example1.txt example2.txt | sort -nk1
//
// will output
//
//	1 also
//	2 an
//	2 example
//	2 is
//	2 this
//
// For information on which options are available, run
//
//	occurs --help
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

	seq = flag.Bool("seq", false, "Count files sequentially, rather than in parallel.")

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

	if *seq {
		countSeq()
	} else {
		countParallel()
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
