package application

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"
)

type UsecaseGetStockValue struct {
	stockQuote domain.StockQuote
	stockEvent domain.StockEvent
}

func NewUsecaseGetStockValue(s domain.StockQuote, e domain.StockEvent) UsecaseGetStockValue {
	return UsecaseGetStockValue{
		stockQuote: s,
		stockEvent: e,
	}
}

func (u UsecaseGetStockValue) Execute(ctx context.Context, input UsecaseGetStockValueInput) (domain.Stock, error) {
	stock, err := u.stockQuote.Get(ctx, domain.NewStock(input.Code, ""))
	if err != nil {
		return domain.Stock{}, err
	}

	err = u.stockEvent.StockRequested(ctx, stock)
	if err != nil {
		return domain.Stock{}, err
	}

	return stock, nil
}
