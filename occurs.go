package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {
	tl := flag.Bool("tl", false, "Trim whitespace from lines.")
	ic := flag.Bool("ic", false, "Ignore case. (All lines will be converted to lower case.)")
	se := flag.Bool("se", false, "Ignore empty lines.")
	flag.Parse()

	filters := make(Filters, 0, 2)
	if *tl {
		filters = append(filters, FilterFunc(strings.TrimSpace))
	}
	if *ic {
		filters = append(filters, FilterFunc(strings.ToLower))
	}

	c := &Counter{
		Filters:   filters,
		SkipEmpty: *se,
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

	for _, arg := range args {
		if arg == "-" {
			c.Count(os.Stdin)
			continue
		}

		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping %q because of error: %v\n", arg, err)
			continue
		}

		c.Count(file)

		file.Close()
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', tabwriter.StripEscape)
	for line, num := range c.Lines {
		fmt.Fprintf(tw, "%v\t\xFF%v\xFF\n", num, line)
	}
	tw.Flush()
}
