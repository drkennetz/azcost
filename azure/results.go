package azure

import (
	"strings"
)

// CostResult holds the result of
type RawCostResult struct {
	Cost     float64
	Date     string
	Resource string
	Currency string
}

type CostResults struct {
	Resources []RawCostResult
}

type ParsedCostResult struct {
	Cost           float64
	Date           string
	ParsedResource string
}

type ParsedCostResults struct {
	Results []ParsedCostResult
}

const dateFormat = "2006-01-02"

// ParseIdResults parses the results of resourceId query and begins grouping
func (results *CostResults) ParseIdResults() ParsedCostResults {
	var parsedResults ParsedCostResults
	for _, v := range results.Resources {
		var parsedResult ParsedCostResult
		parsedResult.Cost = v.Cost
		parsedResult.Date = v.Date
		splitResource := strings.Split(v.Resource, "/")
		resourceType := splitResource[len(splitResource)-2]
		resourceName := strings.Split(splitResource[len(splitResource)-1], "-")[0]
		parsedResult.ParsedResource = strings.Join([]string{resourceType, resourceName}, "/")
		parsedResults.Results = append(parsedResults.Results, parsedResult)
	}
	return parsedResults
}
