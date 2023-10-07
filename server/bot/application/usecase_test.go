package application_test

import (
	"context"
	"errors"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg/bot/mock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUsecaseGetStockValue(t *testing.T) {
	tests := []struct {
		name                 string
		stockQuoteDependency func(ctrl *gomock.Controller) domain.StockQuote
		input                application.UsecaseGetStockValueInput
		stockExpected        domain.Stock
		errExpected          error
	}{
		{
			name: "Test get stock value with success",
			stockQuoteDependency: func(ctrl *gomock.Controller) domain.StockQuote {
				service := mock.NewMockStockQuote(ctrl)
				service.
					EXPECT().
					Get(
						gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
						gomock.AssignableToTypeOf(domain.Stock{}),
					).
					Return(domain.NewStock("APPL.US", "$93.42"), nil).
					Times(1)

				return service
			},
			stockExpected: domain.Stock{
				Code:  "APPL.US",
				Value: "$93.42",
			},
			errExpected: nil,
		},
		{
			name: "Test get stock value with unknown error from stock quote on get",
			stockQuoteDependency: func(ctrl *gomock.Controller) domain.StockQuote {
				service := mock.NewMockStockQuote(ctrl)
				service.
					EXPECT().
					Get(
						gomock.AssignableToTypeOf(reflect.TypeOf((*context.Context)(nil)).Elem()),
						gomock.AssignableToTypeOf(domain.Stock{}),
					).
					Return(domain.Stock{}, errors.New("failed to get stock data")).
					Times(1)

				return service
			},
			stockExpected: domain.Stock{},
			errExpected:   errors.New("failed to get stock data"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			ctrl := gomock.NewController(t)

			usecase := application.NewUsecaseGetStockValue(test.stockQuoteDependency(ctrl))
			stockGot, errGot := usecase.Execute(ctx, test.input)

			assert.Equal(t, test.stockExpected, stockGot)
			assert.Equal(t, test.errExpected, errGot)
		})
	}
}
