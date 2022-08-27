package azure

import (
	"context"
	"flag"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
	"github.com/Azure/go-autorest/autorest/to"
	"log"
	"strconv"
	"time"
)

var (
	subscription = flag.String("subscriptionid", "", "subscription id, format: 00000000-0000-0000-0000-000000000000")
	start        = flag.String("start", "", "start date of range to measure cost")
	end          = flag.String("end", "", "end date of range to measure cost")
)

// RunType returns cost results by either ResourceType of ResourceId
func RunType(name string) CostResults {
	flag.Parse()
	if *subscription == "" {
		log.Fatalln("A valid azure subscription which you can access is required ")
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalln(err)
	}
	if *start == "" {
		// set to yesterday
		*start = time.Now().AddDate(0, 0, -1).Format(dateFormat)
	}
	if *end == "" {
		// set to today
		*end = time.Now().Format(dateFormat)
	}

	begin, err := time.Parse(dateFormat, *start)
	if err != nil {
		log.Fatalln(err)
	}
	stop, err := time.Parse(dateFormat, *end)
	if err != nil {
		log.Fatalln(err)
	}
	costClient, err := armcostmanagement.NewQueryClient(cred, nil)
	if err != nil {
		log.Fatalln(err)
	}

	aggregation := make(map[string]*armcostmanagement.QueryAggregation)
	sum := armcostmanagement.QueryAggregation{
		Function: (*armcostmanagement.FunctionType)(to.StringPtr("Sum")),
		Name:     to.StringPtr("Cost"),
	}
	aggregation["totalCost"] = &sum

	var grouping []*armcostmanagement.QueryGrouping
	grouping = append(grouping, &armcostmanagement.QueryGrouping{
		Name: to.StringPtr(name),
		Type: (*armcostmanagement.QueryColumnType)(to.StringPtr("Dimension")),
	})
	queryDefinition := armcostmanagement.QueryDefinition{
		Type:      (*armcostmanagement.ExportType)(to.StringPtr("ActualCost")),
		Timeframe: (*armcostmanagement.TimeframeType)(to.StringPtr("Custom")),
		Dataset: &armcostmanagement.QueryDataset{
			Granularity: (*armcostmanagement.GranularityType)(to.StringPtr("daily")),
			Aggregation: aggregation,
			Grouping:    grouping,
		},
		TimePeriod: &armcostmanagement.QueryTimePeriod{
			From: &begin,
			To:   &stop,
		},
	}

	subscriptionId := "/subscriptions/" + *subscription
	results, err := costClient.Usage(context.Background(), subscriptionId, queryDefinition, nil)
	if err != nil {
		log.Fatalln(err)
	}

	var costResults CostResults
	// Parse data
	for _, v := range results.Properties.Rows {
		var result CostResult
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
					panic("Unknown type")
				}
			case 2:
				fl, _ := v2.(string)
				result.Resource = fl
			case 3:
				fl, _ := v2.(string)
				result.Currency = fl
			}
		}
		costResults.Resources = append(costResults.Resources, result)
	}
	return costResults
}
