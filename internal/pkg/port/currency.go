package port

import (
	"github.com/darcops/boletia/internal/pkg/domain"
)

// CurrencyService is the signature port to perform operations over the currency resource.
type CurrencyService interface {
	GetHistory(currency, startDate, endDate string) ([]*domain.CurrenciesHistory, error)
}

// CurrencyRepository is the signature port to perform data base operations.
type CurrencyRepository interface {
	Find(out interface{}, conditions ...interface{}) error
	CreateInBatches(value interface{}, batchSize int) error
	GetMinAndMaxDatesInHistory() (int64, int64)
}
