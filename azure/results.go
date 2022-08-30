package azure

import (
	"fmt"
	"strings"
)

type RawCostResultNoId struct {
	Cost          float64
	Date          string
	ResourceType  string
	ResourceGroup string
	Currency      string
}

type CostResultsNoId struct {
	Resources []RawCostResultNoId
}

type ParsedCostResultNoId struct {
	Cost                float64
	Date                string
	ParsedResourceType  string
	ParsedResourceGroup string
}

type ParsedCostResultsNoId struct {
	Results []ParsedCostResultNoId
}

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

// ParseIdResults parses the results of resourceId query and begins grouping
func (results *CostResultsNoId) ParseNoIdResults() ParsedCostResultsNoId {
	var parsedResults ParsedCostResultsNoId
	for _, v := range results.Resources {
		var parsedResult ParsedCostResultNoId
		parsedResult.Cost = v.Cost
		parsedResult.Date = v.Date
		parsedResult.ParsedResourceGroup = v.ResourceGroup
		parsedResult.ParsedResourceType = v.ResourceType

		if parsedResult.ParsedResourceGroup == "" {
			parsedResult.ParsedResourceGroup = "MicrosoftInternal"
		}
		if parsedResult.ParsedResourceType == "" {
			parsedResult.ParsedResourceType = "reservations"
		}
		parsedResults.Results = append(parsedResults.Results, parsedResult)
	}
	return parsedResults
}

type GroupBy struct {
	Gb map[string]map[string]float64
}

func (results *ParsedCostResultsNoId) GroupByRg() GroupBy {
	var gb GroupBy
	gb.Gb = make(map[string]map[string]float64)
	var totalCost float64
	for _, v := range results.Results {
		if _, ok := gb.Gb[v.ParsedResourceGroup]; ok {
			gb.Gb[v.ParsedResourceGroup][v.ParsedResourceType] += v.Cost
			totalCost += v.Cost
		} else {
			gb.Gb[v.ParsedResourceGroup] = make(map[string]float64)
			gb.Gb[v.ParsedResourceGroup][v.ParsedResourceType] += v.Cost
			totalCost += v.Cost
		}
	}
	fmt.Println("total cost in gb: ", totalCost)
	return gb
}

func (gb *GroupBy) PrettyPrint() {
	for key := range gb.Gb {
		fmt.Println(key, ":")
		for subkey := range gb.Gb[key] {
			fmt.Printf("\t%v: %v\n", subkey, gb.Gb[key][subkey])
		}
	}
}
