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
	subscription = flag.String("subscriptionid", "", "subscription id, format: 00000000-0000-0000-0000-000000000000")
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

	if *subscription == "" {
		log.Fatalln("Please include a valid azure subscription id of the form: 00000000-0000-0000-0000-000000000000 ")
	}
	if *start == "" {
		log.Fatalln("Please include start date with format YYYY-MM-DD")
	}
	if *end == "" {
		log.Fatalln("Please include end date with format YYYY-MM-DD")
	}
	if *r1 {
		grouping := azure.NewResourceIdTypeGroupGrouping("Dimension")
		results := azure.Run(*start, *end, *subscription, grouping)
		parsedResults := results.ParseIdResults()
		for _, v := range parsedResults.Results {
			fmt.Println(v.Date, v.Cost, v.ParsedResourceGroup, v.ParsedResourceType, v.ParsedResourceId)
		}
	}
	if *r2 {
		grouping := azure.NewResourceTypeGroupGrouping("Dimension")
		results := azure.Run2(*start, *end, *subscription, grouping)
		var totalCostResource float64
		for _, v := range results.Resources {
			totalCostResource += v.Cost
		}
		parsedResults := results.ParseNoIdResults()
		gb := parsedResults.GroupByRg()
		gb.PrettyPrint()
		fmt.Println("Total cost of resources for time period", *start, *end, ":", totalCostResource)
	}
	os.Exit(0)
}
