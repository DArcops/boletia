package postgres

import (
	"fmt"
	"sync"
	"time"

	"github.com/darcops/boletia/internal/pkg/config"
	"github.com/darcops/boletia/internal/pkg/domain"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var onceDBLoad sync.Once

var tables = []interface{}{
	&domain.CurrenciesHistory{},
}

func connect() *gorm.DB {
	onceDBLoad.Do(func() {
		source := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			config.CfgIn.PostgresHost,
			config.CfgIn.PostgresUser,
			config.CfgIn.PostgresPass,
			config.CfgIn.PostgresName,
			config.CfgIn.PostgresPort,
		)
		var i int
		for {
			var err error
			if i >= 30 {
				panic("could not connect to PostgreSQL" + source)
			}
			time.Sleep(3 * time.Second)

			db, err = gorm.Open(postgres.Open(source), &gorm.Config{})

			if err != nil {
				log.Info("Retrying connection...", err)
				i++
				continue
			}

			break
		}
		migrate()
		log.Info("Connected to db!")
	})

	return db
}

func migrate() {
	for _, table := range tables {
		db.AutoMigrate(table)
	}
}
