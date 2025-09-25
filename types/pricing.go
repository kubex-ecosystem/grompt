package types

// Pricing represents the pricing details of a provider.
type Pricing struct {
	CostPerThousandTokens float64 `json:"cost_per_thousand_tokens"`
	InputCostPer1K        float64 `json:"input_cost_per_1k"`
	OutputCostPer1K       float64 `json:"output_cost_per_1k"`
	Currency              string  `json:"currency"`
}
