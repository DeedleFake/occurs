package main

type Filterer interface {
	Filter(string) string
}

type FilterFunc func(string) string

func (ff FilterFunc) Filter(str string) string {
	return ff(str)
}

type Filters []Filterer

func (f Filters) Filter(str string) string {
	for _, filter := range f {
		str = filter.Filter(str)
	}

	return str
}
