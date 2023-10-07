package application

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"
)

type UsecaseGetStockValue struct {
	stockQuote domain.StockQuote
}

func NewUsecaseGetStockValue(s domain.StockQuote) UsecaseGetStockValue {
	return UsecaseGetStockValue{
		stockQuote: s,
	}
}

func (u UsecaseGetStockValue) Execute(ctx context.Context, input UsecaseGetStockValueInput) (domain.Stock, error) {
	stock, err := u.stockQuote.Get(ctx, domain.NewStock(input.Code, ""))
	if err != nil {
		return domain.Stock{}, err
	}

	return stock, nil
}
