package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/drkennetz/azcost/azure"
	"os"
)

var (
	printVersion = flag.Bool("version", false, "print version and exit")
)

//go:embed VERSION
var version string

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	results := azure.RunType()
	for _, v := range results {
		fmt.Println(v)
	}
}
