package azure

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement"

// NewResourceIdTypeGroupGrouping returns a new QueryGrouping querying ResourceId, ResourceType, and ResourceGroup
func NewResourceIdTypeGroupGrouping(queryColumnType string) []*armcostmanagement.QueryGrouping {
	var grouping []*armcostmanagement.QueryGrouping
	newGroupingId := NewQueryGrouping(QueryGroupingResourceId, queryColumnType)
	newGroupingType := NewQueryGrouping(QueryGroupingResourceType, queryColumnType)
	newGroupingGroup := NewQueryGrouping(QueryGroupingResourceGroup, queryColumnType)
	grouping = append(grouping, &newGroupingId, &newGroupingType, &newGroupingGroup)
	return grouping
}

// NewResourceTypeGroupGrouping returns a new QueryGrouping querying ResourceType and ResourceGroup
func NewResourceTypeGroupGrouping(queryColumnType string) []*armcostmanagement.QueryGrouping {
	var grouping []*armcostmanagement.QueryGrouping
	newGroupingType := NewQueryGrouping(QueryGroupingResourceType, queryColumnType)
	newGroupingGroup := NewQueryGrouping(QueryGroupingResourceGroup, queryColumnType)
	grouping = append(grouping, &newGroupingType, &newGroupingGroup)
	return grouping
}
