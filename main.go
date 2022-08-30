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
	resourceId   = flag.Bool("i", false, "print cost by resource id")
	resourceType = flag.Bool("r", false, "print cost by resource type")
)

//go:embed VERSION
var version string

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if *resourceId {
		results := azure.Run("ResourceId")
		parsedResults := results.ParseIdResults()
		for _, v := range parsedResults.Results {
			fmt.Println(v.Date, v.ParsedResource, v.Cost)
		}
	}
	if *resourceType {
		results := azure.Run("ResourceType")
		for _, v := range results.Resources {
			fmt.Println(v)
		}
	}
	os.Exit(0)
}
