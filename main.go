package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/drkennetz/azcost/azure"
	"github.com/drkennetz/azcost/utils"
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
	fn           = flag.String("f", "", "filename of results - if not included will be <start>.<end>.csv")
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
	if *fn == "" {
		*fn = fmt.Sprintf("%s.%s.csv", *start, *end)
	}
	if *r1 {
		var nextLink string
		var allResults azure.ParsedCostResults
		grouping := azure.NewResourceIdTypeGroupGrouping("Dimension")
		initResults, nextLink := azure.Run(*start, *end, *subscription, nextLink, grouping)
		parsedResults := initResults.ParseIdResults()
		allResults.Results = append(allResults.Results, parsedResults.Results...)
		for nextLink != "" {
			tmpResults, tmpLink := azure.Run(*start, *end, *subscription, nextLink, grouping)
			parsedResults = tmpResults.ParseIdResults()
			allResults.Results = append(allResults.Results, parsedResults.Results...)
			nextLink = tmpLink
		}
		var totalCost float64
		for _, v := range allResults.Results {
			totalCost += v.Cost
		}
		fmt.Println("total cost: ", totalCost)
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
		utils.WriteCSV(*fn, gb)
	}
	os.Exit(0)
}
