package postgres

import (
	"sync"

	"gorm.io/gorm"
)

// ClientDB represents a postgres client
type clientDB struct {
	db    *gorm.DB
	write *gorm.DB
}

var (
	oncePostgresClient sync.Once
	postgresClient     *clientDB
)

// NewRepository returns a new instance for postgres db
func NewRepository() *clientDB {
	oncePostgresClient.Do(func() {
		postgresClient = &clientDB{
			db: connect(),
		}
	})
	return postgresClient
}

// Find loads results from db in the first parameter
// Accoirding the given filters
func (c *clientDB) Find(out interface{}, where ...interface{}) error {
	return c.db.Debug().Find(out, where...).Error
}

func (c *clientDB) Create(value interface{}) error {
	return c.db.Create(value).Error
}

func (c *clientDB) CreateInBatches(value interface{}, batchSize int) error {
	return db.CreateInBatches(value, batchSize).Error
}

func (c *clientDB) GetMinAndMaxDatesInHistory() (int64, int64) {
	var minVal, maxVal int64

	db.Table("currencies.currencies_history").Select("MIN(timestamp) as min_val").Scan(&minVal)
	db.Table("currencies.currencies_history").Select("MAX(timestamp) as max_val").Scan(&maxVal)

	return minVal, maxVal
}

func (c *clientDB) GetAllCurrencies(value interface{}, start, end int64) error {
	return db.Find(value, "timestamp BETWEEN ? AND ?", start, end).Error
}

func (c *clientDB) GetCurrencyHistory(value interface{}, currency string, start, end int64) error {
	return db.Find(value, "timestamp BETWEEN ? AND ? AND currency_code = ?", start, end, currency).Error
}
