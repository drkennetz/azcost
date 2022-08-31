package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
	"log"
	"strconv"
	"time"
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

// Run2 returns cost results for a given time range
func Run2(start, end, subscriptionid string, grouping []*armcostmanagement.QueryGrouping) CostResultsNoId {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalln(err)
	}

	begin, err := time.Parse(dateFormat, start)
	if err != nil {
		log.Fatalln(err)
	}
	stop, err := time.Parse(dateFormat, end)
	if err != nil {
		log.Fatalln(err)
	}
	costClient, err := armcostmanagement.NewQueryClient(cred, nil)
	if err != nil {
		log.Fatalln(err)
	}

	aggregation := make(map[string]*armcostmanagement.QueryAggregation)
	sum := NewQueryAggregation("Sum", "Cost")
	aggregation["totalCost"] = &sum
	newQueryDefinition := NewQueryDefinition("ActualCost", "Custom", "daily", aggregation, grouping, begin, stop)
	subscriptionId := "/subscriptions/" + subscriptionid
	results, err := costClient.Usage(context.Background(), subscriptionId, newQueryDefinition, nil)
	if err != nil {
		log.Fatalln(err)
	}
	var costResults CostResultsNoId
	// Parse data
	for _, v := range results.Properties.Rows {
		var result RawCostResultNoId
		for i, v2 := range v {
			switch i {
			case 0:
				fl, _ := v2.(float64)
				result.Cost = fl
			case 1:
				if str, ok := v2.(string); ok {
					date, err := time.Parse("2006-01-02T15:04:05", str)
					if err != nil {
						log.Fatalln(err)
					}
					result.Date = date.Format(dateFormat)
				} else if str, ok := v2.(float64); ok {
					var y int64 = int64(str)
					timestring := strconv.Itoa(int(y))
					time, err := time.Parse("20060102", timestring)
					if err != nil {
						log.Fatalln(err)
					}
					result.Date = time.Format(dateFormat)

				} else {
					log.Fatalln("Unknown type")
				}
			case 2:
				fl, _ := v2.(string)
				result.ResourceType = fl
			case 3:
				fl, _ := v2.(string)
				result.ResourceGroup = fl
			case 4:
				fl, _ := v2.(string)
				result.Currency = fl
			}
		}
		costResults.Resources = append(costResults.Resources, result)
	}
	return costResults
}
