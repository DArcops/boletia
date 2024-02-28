package currency

import (
	"net/http"

	"github.com/darcops/boletia/internal/pkg/port"
	er "github.com/darcops/errors"
	"github.com/gin-gonic/gin"
)

type currencyHandler struct {
	service port.CurrencyService
}

func newHandler(service port.CurrencyService) *currencyHandler {
	return &currencyHandler{
		service: service,
	}
}

func (c *currencyHandler) getCurrencies(ctx *gin.Context) {
	currency := ctx.Param("currencyCode")
	finit := ctx.Query("finit")
	fend := ctx.Query("fend")

	currencies, err := c.service.GetHistory(currency, finit, fend)
	if err != nil {
		er.JSON(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}
