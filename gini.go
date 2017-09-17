package gini

import ()

type Gini struct {
	filepath string
}

func ParseFile(filepath string) *Gini {
	gini := &Gini{
		filepath: filepath,
	}
	gini.Parse()

	return gini
}

func Parse(src string) string {
	return src
}

