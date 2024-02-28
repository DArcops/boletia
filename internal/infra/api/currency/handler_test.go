package currency

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darcops/boletia/internal/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock de CurrencyService para pruebas
type mockCurrencyService struct{}

func (m *mockCurrencyService) GetHistory(currency, startDate, endDate string) ([]*domain.CurrenciesHistory, error) {
	if currency == "USD" {
		return []*domain.CurrenciesHistory{}, nil
	}
	return nil, errors.New("Currency not found")
}

func TestGetCurrencies(t *testing.T) {
	// Configurar el enrutador Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Crear un manejador de moneda con un servicio de moneda ficticio
	handler := newHandler(&mockCurrencyService{})

	// Ruta de prueba
	router.GET("/currencies/:currencyCode", handler.getCurrencies)

	// Caso de prueba 1: Moneda válida
	w := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/currencies/USD", nil)
	router.ServeHTTP(w, req1)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200")
}
