package main

import "encoding/json"

const exchangeRateUrl = "https://www.tinkoff.ru/api/trading/currency/list?sessionId="

type ExchangeRateRequest struct {
	Country     string `json:"country"`
	CurrentPage int    `json:"currentPage"`
	End       int    `json:"end"`
	OrderType string `json:"orderType"`
	PageSize int    `json:"pageSize"`
	SortType string `json:"sortType"`
	Start    int    `json:"start"`
}

type ExchangeRateResponsePayloadValue struct {
	Price struct {
		Currency string  `json:"currency"`
		Value    float64 `json:"value"`
		FromCache bool    `json:"from_cache"`
	} `json:"price"`
	HistoricalPrices []struct {
		Amount   float64 `json:"amount"`
		UnixTime float64 `json:"unixtime"`
	} `json:"historicalPrices"`
	Symbol struct {
		Ticker string `json:"ticker"`
	} `json:"symbol"`
}

type ExchangeRateResponse struct {
	Payload struct {
		Total int `json:"total"`
		Values []ExchangeRateResponsePayloadValue `json:"values"`
	} `json:"payload"`
	Status string `json:"status"`
}

func getExchangeRate(token string) map[string]float64 {
	request := ExchangeRateRequest{
		Country: "All",
		CurrentPage: 0,
		End: 12,
		OrderType: "Asc",
		PageSize: 12,
		SortType: "ByBuyBackDate",
		Start: 0,
	}

	exchangeRateParcel := ExchangeRateResponse{}

	requestBytes, _ := json.Marshal(request)

	postJson(exchangeRateUrl + token, &exchangeRateParcel, requestBytes)

	var result = map[string]float64{}

	for i := 0; i < len(exchangeRateParcel.Payload.Values); i++ {
		var val = exchangeRateParcel.Payload.Values[i]
		result[val.Symbol.Ticker] = val.HistoricalPrices[12].Amount
	}

	return result
}

