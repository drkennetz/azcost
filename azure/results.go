package azure

import (
	"strings"
)

// CostResult holds the result of
type RawCostResult struct {
	Cost          float64
	Date          string
	ResourceId    string
	ResourceType  string
	ResourceGroup string
	Currency      string
}

type CostResults struct {
	Resources []RawCostResult
}

type ParsedCostResult struct {
	Cost                float64
	Date                string
	ParsedResourceId    string
	ParsedResourceType  string
	ParsedResourceGroup string
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
		parsedResult.ParsedResourceGroup = v.ResourceGroup
		parsedResult.ParsedResourceType = v.ResourceType
		parsedResult.ParsedResourceId = v.ResourceId

		if parsedResult.ParsedResourceGroup == "" {
			parsedResult.ParsedResourceGroup = "MicrosoftInternal"
		}
		if parsedResult.ParsedResourceType == "" {
			parsedResult.ParsedResourceType = "reservations"
		}
		tmpId := strings.Split(parsedResult.ParsedResourceId, "/")
		lastIndex := len(tmpId) - 1
		if tmpId[lastIndex] == "" {
			parsedResult.ParsedResourceId = tmpId[lastIndex-1]
		} else {
			parsedResult.ParsedResourceId = tmpId[lastIndex]
		}
		parsedResults.Results = append(parsedResults.Results, parsedResult)
	}
	return parsedResults
}
