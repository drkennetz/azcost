package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/drkennetz/azcost/azure"
	"log"
	"os"
)

var (
	printVersion = flag.Bool("version", false, "print version and exit")
	start        = flag.String("start", "", "start date of range to measure cost")
	end          = flag.String("end", "", "end date of range to measure cost")
)

//go:embed VERSION
var version string

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *start == "" {
		log.Fatalln("Please include start date with format YYYY-MM-DD")
	}
	if *end == "" {
		log.Fatalln("Please include end date with format YYYY-MM-DD")
	}
	results := azure.Run(*start, *end)
	parsedResults := results.ParseIdResults()
	for _, v := range parsedResults.Results {
		fmt.Println(v.Date, v.Cost, v.ParsedResourceGroup, v.ParsedResourceType, v.ParsedResourceId)
	}
	os.Exit(0)
}
