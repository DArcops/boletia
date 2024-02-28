package currency

import (
	"testing"

	"github.com/darcops/boletia/internal/pkg/domain"
	"github.com/darcops/boletia/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type currencyServiceSuite struct {
	// we need this to use the suite functionalities from testify
	suite.Suite

	// the generated mocked version of our repo
	repo *mocks.CurrencyRepository

	// the functionalities we want to test
	service *service
}

// SetupTest runs before any test
func (suite *currencyServiceSuite) SetupTest() {
	// init the mocked version of repo
	repo := new(mocks.CurrencyRepository)

	// inject repo to the service, since service needs repo to work
	service := NewService(repo)

	// assign them as suit properties
	suite.repo = repo
	suite.service = service
}

func (suite *currencyServiceSuite) TestGetHistory_All() {
	var expectedResponse []*domain.CurrenciesHistory

	mockedCurrencies := []*domain.CurrenciesHistory{
		{CurrencyCode: "MXN", Value: 20.0, Timestamp: 12345},
		{CurrencyCode: "GMD", Value: 54.0, Timestamp: 12345},
		{CurrencyCode: "CUP", Value: 26.0, Timestamp: 12345},
	}

	suite.repo.On("GetMinAndMaxDatesInHistory").Return(int64(12345), int64(12345)).Once()

	suite.repo.On("Find",
		&expectedResponse,
		mock.Anything,
		mock.Anything,
		mock.Anything).Run(func(args mock.Arguments) {
		r := args.Get(0).(*[]*domain.CurrenciesHistory)
		*r = mockedCurrencies
		expectedResponse = *r
	}).Return(nil)

	// Calling method that we want to test.
	got, err := suite.service.GetHistory("all", "", "")

	// Assertions
	suite.Nil(err)
	suite.ElementsMatch(expectedResponse, got)
}

func (suite *currencyServiceSuite) TestGetHistory_With_Some_Currency() {
	var expectedResponse []*domain.CurrenciesHistory

	mockedCurrencies := []*domain.CurrenciesHistory{
		{CurrencyCode: "MXN", Value: 20.0, Timestamp: 12345},
		{CurrencyCode: "MXN", Value: 26.1, Timestamp: 12345},
	}

	suite.repo.On("GetMinAndMaxDatesInHistory").Return(int64(12345), int64(12345)).Once()

	suite.repo.On("Find",
		&expectedResponse,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).Run(func(args mock.Arguments) {
		r := args.Get(0).(*[]*domain.CurrenciesHistory)
		*r = mockedCurrencies
		expectedResponse = *r
	}).Return(nil)

	// Calling method that we want to test.
	got, err := suite.service.GetHistory("mxn", "", "")

	// Assertions
	suite.Nil(err)
	suite.Equal(len(got), 2)
}

func (suite *currencyServiceSuite) TestGetCurrencies() {
	httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles",
		httpmock.NewStringResponder(200, `{"data": {"MXN": {"code": "MXN", "value": 17.2}}}`))

	suite.repo.On("CreateInBatches", mock.Anything, mock.Anything).Return(nil).Once()

	// Calling method that we want to test.
	suite.service.GetCurrencies()
}

func (suite *currencyServiceSuite) TestParseDate() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"2022-01-01T12:00:00", 1641038400}, // Expected Unix timestamp for "2022-01-01T12:00:00"
		{"2021-12-31T23:59:59", 1640995199}, // Expected Unix timestamp for "2021-12-31T23:59:59"
		{"invalid-date", 0},                 // Expected return value for an invalid date
	}

	for _, test := range tests {
		suite.T().Run(test.input, func(t *testing.T) {
			result := parseDate(test.input)
			if result != test.expected {
				t.Errorf("parseDate(%s) = %d; expected %d", test.input, result, test.expected)
			}
		})
	}

}

func TestClientService(t *testing.T) {
	suite.Run(t, new(currencyServiceSuite))
}
