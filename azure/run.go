package azure

import (
	"context"
	"flag"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
	"log"
	"strconv"
	"time"
)

var (
	subscription = flag.String("subscriptionid", "", "subscription id, format: 00000000-0000-0000-0000-000000000000")
)

// Run returns cost results for a given time range
func Run(start, end string) CostResults {
	flag.Parse()
	if *subscription == "" {
		log.Fatalln("A valid azure subscription which you can access is required ")
	}
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
	var grouping []*armcostmanagement.QueryGrouping
	newGroupingId := NewQueryGrouping(QueryGroupingResourceId, "Dimension")
	newGroupingType := NewQueryGrouping(QueryGroupingResourceType, "Dimension")
	NewGroupingGroup := NewQueryGrouping(QueryGroupingResourceGroup, "Dimension")
	grouping = append(grouping, &newGroupingId, &newGroupingType, &NewGroupingGroup)
	newQueryDefinition := NewQueryDefinition("ActualCost", "Custom", "daily", aggregation, grouping, begin, stop)
	subscriptionId := "/subscriptions/" + *subscription
	results, err := costClient.Usage(context.Background(), subscriptionId, newQueryDefinition, nil)
	if err != nil {
		log.Fatalln(err)
	}
	var costResults CostResults
	// Parse data
	for _, v := range results.Properties.Rows {
		var result RawCostResult
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
				result.ResourceId = fl
			case 3:
				fl, _ := v2.(string)
				result.ResourceType = fl
			case 4:
				fl, _ := v2.(string)
				result.ResourceGroup = fl
			case 5:
				fl, _ := v2.(string)
				result.Currency = fl
			}
		}
		costResults.Resources = append(costResults.Resources, result)
	}
	return costResults
}
