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
	r1           = flag.Bool("r1", false, "aggregate by resourceId")
	r2           = flag.Bool("r2", false, "aggregate by resourceType")
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
	if *r1 {
		results := azure.Run(*start, *end)
		parsedResults := results.ParseIdResults()
		for _, v := range parsedResults.Results {
			fmt.Println(v.Date, v.Cost, v.ParsedResourceGroup, v.ParsedResourceType, v.ParsedResourceId)
		}
	}
	if *r2 {
		results := azure.Run2(*start, *end)
		var totalCostResource float64
		for _, v := range results.Resources {
			totalCostResource += v.Cost
		}
		var totalCostParsed float64
		parsedResults := results.ParseNoIdResults()
		for _, v := range parsedResults.Results {
			totalCostParsed += v.Cost
		}
		gb := parsedResults.GroupByRg()
		gb.PrettyPrint()
		fmt.Println("total cost resource: ", totalCostResource)
		fmt.Println("total cost parsed: ", totalCostParsed)
	}
	os.Exit(0)
}
