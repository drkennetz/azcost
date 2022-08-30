package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"
	"github.com/Azure/go-autorest/autorest/to"
	"time"
)

const (
	QueryGroupingResourceGroup          = "ResourceGroup"
	QueryGroupingResourceGroupName      = "ResourceGroupName"
	QueryGroupingResourceLocation       = "ResourceLocation"
	QueryGroupingConsumedService        = "ConsumedService"
	QueryGroupingResourceType           = "ResourceType"
	QueryGroupingResourceId             = "ResourceId"
	QueryGroupingMeterId                = "MeterId"
	QueryGroupingBillingMonth           = "BillingMonth"
	QueryGroupingMeterCategory          = "MeterCategory"
	QueryGroupingMeterSubcategory       = "MeterSubcategory"
	QueryGroupingMeter                  = "Meter"
	QueryGroupingAccountName            = "AccountName"
	QueryGroupingDepartmentName         = "DepartmentName"
	QueryGroupingSubscriptionId         = "SubscriptionId"
	QueryGroupingSubscriptionName       = "SubscriptionName"
	QueryGroupingServiceName            = "ServiceName"
	QueryGroupingServiceTier            = "ServiceTier"
	QueryGroupingEnrollmentAccountName  = "EnrollmentAccountName"
	QueryGroupingBillingAccountId       = "BillingAccountId"
	QueryGroupingResourceGuid           = "ResourceGuid"
	QueryGroupingBillingPeriod          = "BillingPeriod"
	QueryGroupingInvoiceNumber          = "InvoiceNumber"
	QueryGroupingChargeType             = "ChargeType"
	QueryGroupingPublisherType          = "PublisherType"
	QueryGroupingReservationId          = "ReservationId"
	QueryGroupingReservationName        = "ReservationName"
	QueryGroupingFrequency              = "Frequency"
	QueryGroupingPartNumber             = "PartNumber"
	QueryGroupingCostAllocationRuleName = "CostAllocationRuleName"
	QueryGroupingMarkupRuleName         = "MarkupRuleName"
	QueryGroupingPricingModel           = "PricingModel"
	QueryGroupingBenefitId              = "BenefitId"
	QueryGroupingBenefitName            = "BenefitName"
)

// NewQueryAggregation returns an aggregation to collect data in Azure
func NewQueryAggregation(fn, name string) armcostmanagement.QueryAggregation {
	return armcostmanagement.QueryAggregation{
		Function: (*armcostmanagement.FunctionType)(to.StringPtr(fn)),
		Name:     to.StringPtr(name),
	}
}

// NewQueryGrouping returns a grouping for data in Azure
func NewQueryGrouping(name, queryColumnType string) armcostmanagement.QueryGrouping {
	return armcostmanagement.QueryGrouping{
		Name: to.StringPtr(name),
		Type: (*armcostmanagement.QueryColumnType)(to.StringPtr(queryColumnType)),
	}
}

// NewQueryDefinition returns a definition for query in Azure
func NewQueryDefinition(exportType, timeframeType, granularity string,
	aggregation map[string]*armcostmanagement.QueryAggregation,
	grouping []*armcostmanagement.QueryGrouping,
	begin, stop time.Time) armcostmanagement.QueryDefinition {

	return armcostmanagement.QueryDefinition{
		Type:      (*armcostmanagement.ExportType)(to.StringPtr(exportType)),
		Timeframe: (*armcostmanagement.TimeframeType)(to.StringPtr(timeframeType)),
		Dataset: &armcostmanagement.QueryDataset{
			Granularity: (*armcostmanagement.GranularityType)(to.StringPtr(granularity)),
			Aggregation: aggregation,
			Grouping:    grouping,
		},
		TimePeriod: &armcostmanagement.QueryTimePeriod{
			From: &begin,
			To:   &stop,
		},
	}
}
