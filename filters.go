package main

// A Filterer wraps the Filter method, which is used to modify lines
// as they are read before counting them.
type Filterer interface {
	Filter(string) string
}

// FilterFunc is an adapter that allows a function that matches its
// signature to work as a Filterer.
type FilterFunc func(string) string

func (ff FilterFunc) Filter(str string) string {
	return ff(str)
}

// Filter is a list of Filterers that is itself a Filterer whose
// Filter() method runs all of the Filterer's Filter() methods in
// sequence, returning the result.
type Filters []Filterer

func (f Filters) Filter(str string) string {
	for _, filter := range f {
		str = filter.Filter(str)
	}

	return str
}
