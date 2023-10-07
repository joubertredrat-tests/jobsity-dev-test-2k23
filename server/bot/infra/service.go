package infra

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"joubertredrat-tests/jobsity-dev-test-2k23/bot/domain"
	"net/http"
	"strings"
)

const (
	URL_FORMAT      = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	CSV_RESULT_LINE = 1
	CSV_VALUE_COLUM = 3
)

type StockQuoteStooq struct {
}

func NewStockQuoteStooq() StockQuoteStooq {
	return StockQuoteStooq{}
}

func (s StockQuoteStooq) Get(ctx context.Context, stock domain.Stock) (domain.Stock, error) {
	requestURL := fmt.Sprintf(URL_FORMAT, strings.ToLower(stock.Code))

	response, err := http.Get(requestURL)
	if err != nil {
		return domain.Stock{}, err
	}
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Stock{}, err
	}

	reader := csv.NewReader(bytes.NewBuffer(resBody))
	records, err := reader.ReadAll()
	if err != nil {
		return domain.Stock{}, err
	}

	value := records[CSV_RESULT_LINE][CSV_VALUE_COLUM]
	if value != "N/D" {
		value = fmt.Sprintf("$%s", value)
	}

	stock.Value = value
	return stock, nil
}
