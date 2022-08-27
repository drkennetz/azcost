package azure

// CostResult holds the result of
type CostResult struct {
	Cost     float64
	Date     string
	Resource string
	Currency string
}

type CostResults struct {
	Resources []CostResult
}

const dateFormat = "2006-01-02"
