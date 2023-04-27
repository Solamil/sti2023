package sti2023

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Rates struct {
	Code     []string `json:"code"`
	Volume   []string `json:"volume"`
	Value    []string `json:"value"`
	CoinCode string   `json:"coin_code"`
	Number   int      `json:"number"`
	Date     string   `json:"date"`
}

var updated time.Time
var cnbUrl string = "https://czk.michalkukla.xyz/json"
var cnbDir string = "cache"
var cnbFilename string = "cnb"
var rates Rates

func GetCurrencySum(amount float64, coinCode string) string {
	var result string = ""
	if amount <= 0 {
		return result
	}
	updateCurrency()
	if !IsExistCode(coinCode) {
		return result
	}
	rate := GetCurrencyRate(coinCode)
	if len(rate) < 2 {
		return result
	}
	value, _ := strconv.ParseFloat(rate[0], 64)
	volume, _ := strconv.ParseFloat(rate[1], 64)

	total := value * (amount / volume)
	result = fmt.Sprintf("%.2f", total)

	return result
}

func GetCurrencyRate(coinCode string) []string {
	var result []string = []string{}
	updateCurrency()
	ReadJsonFile(cnbDir, cnbFilename, &rates)
	var index int = GetIndex(rates.Code, coinCode)
	if index < 0 {
		return result
	}
	result = append(result, rates.Value[index])
	result = append(result, rates.Volume[index])
	return result
}

func GetDate() string {
	updateCurrency()
	ReadJsonFile(cnbDir, cnbFilename, &rates)
	return rates.Date
}

func IsExistCode(coinCode string) bool {
	updateCurrency()
	if !ReadJsonFile(cnbDir, cnbFilename, &rates) {
		return false
	}
	if i := GetIndex(rates.Code, coinCode); i < 0 {
		return false
	}
	return true
}

func GetCoinCodes() []string {
	updateCurrency()
	if !ReadJsonFile(cnbDir, cnbFilename, &rates) {
		return []string{}
	}
	return rates.Code
}

func updateCurrency() {
	if !ReadJsonFile(cnbDir, cnbFilename, &rates) {
		CurrencyRates()
		return
	}
	now := time.Now()
	tUpdate := time.Date(now.Year(), now.Month(), now.Day(), 14, 35+1, 0, 0, now.Location())

	if (now.Before(tUpdate) && now.Day() == updated.Day() && now.Month() == updated.Month() && now.Year() == updated.Year()) ||
		updated.After(tUpdate) {
		// Take it from cache
		return
	}
	CurrencyRates()
}

func CurrencyRates() {
	d := Newrequest(cnbUrl)
	err := json.Unmarshal([]byte(d), &rates)
	if err != nil {
		fmt.Println(err)
		return
	}
	WriteJsonFile(cnbDir, cnbFilename, rates)
	updated = time.Now()
}

func GetIndex(array []string, value string) int {
	var index int = -1
	for i, v := range array {
		if value == v {
			index = i
			break
		}
	}
	return index
}
