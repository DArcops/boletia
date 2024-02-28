package currency

import (
	"github.com/darcops/boletia/internal/infra/repository/postgres"
	"github.com/darcops/boletia/internal/pkg/service/currency"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	repo := postgres.NewRepository()
	handler := newHandler(currency.NewService(repo))

	currencyRoutes := e.Group("/api/v1/currencies")

	{
		currencyRoutes.GET("/:currencyCode", handler.getCurrencies)
	}
}
