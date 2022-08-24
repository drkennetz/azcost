package utils

// DateStartEnd contains start and end dates for a given interval of time "YYYY-MM-DD"
type DateStartEnd struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}
