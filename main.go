package main

import (
	"time"

	"github.com/darcops/boletia/internal/infra/api"
	"github.com/darcops/boletia/internal/infra/repository/postgres"
	"github.com/darcops/boletia/internal/pkg/config"
	"github.com/darcops/boletia/internal/pkg/service/currency"
)

func main() {
	currenciesRepo := postgres.NewRepository()
	currenciesService := currency.NewService(currenciesRepo)

	minutes := time.Duration(config.CfgIn.Minutes)

	go func() {
		for range time.Tick(time.Minute * minutes) {
			currenciesService.GetCurrencies()
		}
	}()

	api.RunServer()
}
