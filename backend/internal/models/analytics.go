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

type NewBuildPremiumResult struct {
	Region         string  `json:"region"`
	NewAvg         int64   `json:"new_avg"`
	OldAvg         int64   `json:"old_avg"`
	PremiumPercent float64 `json:"premium_percent"`
}

type PropertyTypeDistributionResult struct {
	PropertyType string  `json:"property_type"`
	Count        int64   `json:"count"`
	Percentage   float64 `json:"percentage"`
}

type PriceBracketResult struct {
	Bracket    string  `json:"bracket"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type TopActiveAreaResult struct {
	Region           string `json:"region"`
	TransactionCount int64  `json:"transaction_count"`
	TotalValue       int64  `json:"total_value"`
}

type TimeRangeResult struct {
	MinYear int `json:"min_year"`
	MaxYear int `json:"max_year"`
}
