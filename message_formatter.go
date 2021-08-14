package main

import "fmt"

const positiveTrendEmoji = "\U0001F4C8"
const negativeTrendEmoji = "\U0001F4C9"

var signs = map[string]string{
	"USDRUB": "$",
	"EURRUB": "€",
}

func getCurrencyString(ticker string, newValue float64, oldValue float64) string {
	sign := signs[ticker]

	var result = fmt.Sprintf("%s <a href=\"https://www.tinkoff.ru/invest/currencies/%s/\">%.2f</a>", sign, ticker, newValue)

	if oldValue == 0 {
		return result
	}

	var delta = fmt.Sprintf("%.2f", newValue - oldValue)
	if delta == "0.00" {
		return result
	}

	var emoji = positiveTrendEmoji
	if newValue - oldValue < 0 {
		emoji = negativeTrendEmoji
	}

	return fmt.Sprintf("%s %s%s", result, emoji, delta)
}