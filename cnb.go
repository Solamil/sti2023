package sti2023

import (
	"encoding/json"
	"time"
	"strconv"
	"fmt"
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
var url string = "https://czk.michalkukla.xyz/json"
var cnbDir string = "cache"
var cnbFilename string = "cnb"

func GetCurrencySum(amount float64, coinCode string) string {
	var result string = ""

	updateCurrency()

	rate := GetCurrencyRate(coinCode)
	value, err := strconv.ParseFloat(rate[0], 64)
	if err != nil {
		fmt.Println(err)
	}
	volume, err := strconv.ParseFloat(rate[1], 64)
	if err != nil {
		fmt.Println(err)
	}
	total := value * (amount/volume)
	result = fmt.Sprintf("%.2f", total)

	return result
}

func GetCurrencyRate(coinCode string) []string {
	var result []string = []string{}
	updateCurrency()
	var rates Rates
	ReadJsonFile(cnbDir, cnbFilename, &rates)
	var index int = getIndex(rates.Code, coinCode)
	if index < 0 {
		return result
	}
	result = append(result, rates.Value[index])
	result = append(result, rates.Volume[index])
	return result
}

func GetDate() string {
	var rates Rates	
	ReadJsonFile(cnbDir, cnbFilename, &rates)
	return rates.Date
}

func IsExistCode(coinCode string) bool {
	var rates Rates
	if ! ReadJsonFile(cnbDir, cnbFilename, &rates) {
		return false
	}
	if i := getIndex(rates.Code, coinCode); i < 0 {
		return false
	}
	return true
}

func updateCurrency() {
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
	var rates Rates
	d := Newrequest(url)
	err := json.Unmarshal([]byte(d), &rates)
	if err != nil {
		fmt.Println(err)
		return
	}
	WriteJsonFile(cnbDir, cnbFilename, rates)
	updated = time.Now()
}

func getIndex(array []string, value string) int {
	var index int = -1
	for i, v := range array {
		if value == v {
			index = i
			break
		}
	}
	return index
}
