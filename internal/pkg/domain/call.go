package domain

// CurrenciesHistory represents the proper data to be stored in our repository.
// This data type will be the history for the API calls to get currency rates.
type CurrenciesHistory struct {
	ID                    int     `json:"id"`
	CurrencyCode          string  `json:"currencyCode"`
	Value                 float64 `json:"value"`
	Timestamp             int64   `json:"timestamp"`
	LatencyInMilliseconds float64 `json:"latencyInMilliseconds"`
}

// TableName returns a string used by the ORM to set the proper schema and table
// In our repository.
func (CurrenciesHistory) TableName() string {
	return "currencies.currencies_history"
}

type Meta struct {
	LastUpdatedAt string `json:"last_updated_at"`
}

type RateResponseValue struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// ApiCallResponse represents the data returned by the currency API.
type ApiCallResponse struct {
	Meta         `json:"meta"`
	Data         map[string]RateResponseValue `json:"data"`
	EnlapsedTime float64                      `json:"-"`
}
