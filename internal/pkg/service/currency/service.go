package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/darcops/boletia/internal/pkg/config"
	"github.com/darcops/boletia/internal/pkg/domain"
	"github.com/darcops/boletia/internal/pkg/port"
	er "github.com/darcops/errors"
)

const allCurrenciesKey = "ALL"
const batchSize = 100
const maxNumberOfTries = 3

type service struct {
	repository port.CurrencyRepository
}

// NewService returns a new instance of service.
func NewService(repo port.CurrencyRepository) *service {
	return &service{
		repository: repo,
	}
}

// GetHistory returns the history records for the given currency.
// And the given date range.
func (s *service) GetHistory(currency, startDate, endDate string) ([]*domain.CurrenciesHistory, error) {
	var response []*domain.CurrenciesHistory
	var filters []interface{}
	var start, end int64

	currency = strings.ToUpper(currency)

	minDateSaved, maxDateSaved := s.repository.GetMinAndMaxDatesInHistory()

	if startDate == "" || parseDate(startDate) < minDateSaved {
		start = minDateSaved
	} else {
		start = parseDate(startDate)
	}

	if endDate == "" || parseDate(endDate) > maxDateSaved {
		end = maxDateSaved
	} else {
		end = parseDate(endDate)
	}

	if currency != allCurrenciesKey {
		filters = []interface{}{"timestamp BETWEEN ? AND ? AND currency_code = ?", start, end, currency}
	} else {
		filters = []interface{}{"timestamp BETWEEN ? AND ?", start, end}
	}

	return response, s.repository.Find(&response, filters...)
}

func (s *service) GetCurrencies() {
	log.Info("Starting the poolling process...")

	err := s.retryWithBackoff(func() error {
		// Create context with the given timeout
		timeout := time.Duration(config.CfgIn.HttpTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
		defer cancel()

		// Make http call
		return s.doHttpRequest(ctx)
	})

	if err != nil {
		log.Error("Error:", err)
	}

	log.Info("Currencies loaded successfully.")

}

func (s *service) retryWithBackoff(operation func() error) error {
	for attempt := 1; attempt <= maxNumberOfTries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		log.Infof("Error: %v, retrying in %v\n", err, time.Duration(1<<uint(attempt))*time.Second)
		time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
	}
	return fmt.Errorf("Max retry attempts reached")
}

func (s *service) doHttpRequest(ctx context.Context) error {
	url := fmt.Sprintf("%s?apikey=%s", config.CfgIn.ApiHost, config.CfgIn.ApiKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	startTime := time.Now()

	client := http.Client{}
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return er.Build(er.InternalError)
	}
	enlapsedTime := time.Since(startTime).Milliseconds()

	select {
	case <-ctx.Done():
		return fmt.Errorf("Timeout reached: %v", ctx.Err())
	default:
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	result := &domain.ApiCallResponse{}
	if err := json.Unmarshal(body, result); err != nil { // Parse []byte to go struct pointer
		return er.Build(er.InternalError)
	}

	result.EnlapsedTime = float64(enlapsedTime)

	return s.storeResponse(result)
}

func (s *service) storeResponse(res *domain.ApiCallResponse) error {
	var rates []domain.CurrenciesHistory

	for _, rate := range res.Data {
		rates = append(rates, domain.CurrenciesHistory{
			CurrencyCode:          rate.Code,
			Value:                 rate.Value,
			Timestamp:             time.Now().Unix(),
			LatencyInMilliseconds: res.EnlapsedTime,
		})
	}

	return s.repository.CreateInBatches(rates, batchSize)
}

func parseDate(date string) int64 {
	t1, err := time.Parse("2006-01-02T15:04:05", date)
	if err != nil {
		return 0
	}

	return t1.Unix()
}
