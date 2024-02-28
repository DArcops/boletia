package api

import (
	"github.com/darcops/boletia/internal/infra/api/currency"
	"github.com/gin-gonic/gin"
)

func registerRoutes(e *gin.Engine) {
	currency.RegisterRoutes(e)
}
