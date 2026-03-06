package models

type MedianPriceResult struct {
	Region      string `json:"region"`
	MedianPrice int64  `json:"median_price"`
}

type PriceTrendResult struct {
	Period           string `json:"period"`
	AvgPrice         int64  `json:"avg_price"`
	MedianPrice      int64  `json:"median_price"`
	TransactionCount int64  `json:"transaction_count"`
}

type AffordabilityResult struct {
	PropertyType         string  `json:"property_type"`
	AvgPrice             int64   `json:"avg_price"`
	RelativeAffordability float64 `json:"relative_affordability"`
}

type GrowthHotspotResult struct {
	Region        string  `json:"region"`
	GrowthRate    float64 `json:"growth_rate"`
	PrevMedian    int64   `json:"prev_median"`
	CurrentMedian int64   `json:"current_median"`
}
