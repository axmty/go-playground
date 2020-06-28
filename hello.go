package main

import (
	"fmt"

	"github.com/axmty/go-playground/morestrings"
	"github.com/google/go-cmp/cmp"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("Hello world!"))
	fmt.Println(cmp.Diff("Hello world", "Hello go"))
}
