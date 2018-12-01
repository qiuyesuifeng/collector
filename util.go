package main

import (
	"flag"
	"fmt"
	"os"
)

// Usage
func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}
