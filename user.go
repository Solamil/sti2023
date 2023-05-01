package sti2023

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Password  string   `json:"password"`
	Balances  []string `json:"balances"`
	CoinCodes []string `json:"coinCodes"`
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Payment   struct {
		Directions []string `json:"directions"`
		Totals     []string `json:"totals"`
		CoinCodes  []string `json:"coinCodes"`
	} `json:"Payments"`
	SecondAuth struct {
		Code   string `json:"code"`
		Expiry string `json:"expiry"`
	} `json:"secondAuth"`
}

var userDir string = "users"
var user User

func CreatePayment(email string, total float64, direction string, coinCode string) bool {
	var balance float64 = 0.0

	if PreparePayment(email, &balance, &total, &direction, &coinCode) {
		return addPayment(email, balance, total, direction, coinCode)
	}
	return false
}

func PreparePayment(email string, balance, total *float64, direction *string, coinCode *string) bool {
	index := GetIndex(GetUserCoinCodes(email), *coinCode)
	if index < 0 || *total <= 0.0 {
		return false
	}

	balances := GetBalances(email)
	*balance, _ = strconv.ParseFloat(balances[index], 64)

	if strings.EqualFold(*direction, "IN") {
		*balance += *total
	} else if strings.EqualFold(*direction, "OUT") && *balance < *total && !strings.EqualFold(*coinCode, "CZK") {

		valueCzk := GetCurrencySum(*total, *coinCode)
		*total, _ = strconv.ParseFloat(valueCzk, 64)
		index = 0 // CZK
		*coinCode = "CZK"
		*balance, _ = strconv.ParseFloat(balances[index], 64)
		if *balance < *total || *total == 0.0 {
			return false
		}
		*balance -= *total

	} else if strings.EqualFold(*direction, "OUT") && *balance >= *total {
		*balance -= *total
	} else {
		return false
	}
	return true
}

func AddCurrency(email string, code string) bool {
	if !ReadJsonFile(userDir, email, &user) {
		return false
	}
	if index := GetIndex(GetUserCoinCodes(email), code); index > 0 {
		return false
	}
	if !IsExistCode(code) {
		return false
	}
	user.Balances = append(user.Balances, "0.0")
	user.CoinCodes = append(user.CoinCodes, code)
	WriteJsonFile(userDir, email, user)
	return true
}

func GetBalance(email string, code string) string {
	var result string = ""
	var index int = GetIndex(GetUserCoinCodes(email), code)
	if index < 0 {
		return result
	}
	result = GetBalances(email)[index]
	return result
}

func GetBalances(email string) []string {
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.Balances
}

func GetUserCoinCodes(email string) []string {
	var user User
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.CoinCodes
}

func GetNames(email string) []string {
	var user User
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return []string{user.FirstName, user.LastName}
}

func addPayment(email string, balance, total float64, direction, coinCode string) bool {
	if !ReadJsonFile(userDir, email, &user) {
		return false
	}
	index := GetIndex(GetUserCoinCodes(email), coinCode)
	if index < 0 {
		return false
	}
	strTotal := fmt.Sprintf("%.2f", total)
	strBalance := fmt.Sprintf("%.2f", balance)
	user.Balances[index] = strBalance
	user.Payment.Directions = append(user.Payment.Directions, direction)
	user.Payment.Totals = append(user.Payment.Totals, strTotal)
	user.Payment.CoinCodes = append(user.Payment.CoinCodes, coinCode)

	return WriteJsonFile(userDir, email, user)
}

func CheckPassword(email, value string) bool {
	if !ReadJsonFile(userDir, email, &user) {
		return false
	}
	password := fmt.Sprintf("%x", Hash(value))
	return user.Password == password
}

func IsCodeUptodate(email string) bool {
	if !ReadJsonFile(userDir, email, &user) {
		return false
	}

	now := time.Now()
	expiry, err := time.Parse(time.RFC3339, user.SecondAuth.Expiry)
	if err != nil {
		//fmt.Println(err)
		return false
	}
	if expiry.After(now) {
		return true
	}
	return false
}

func CheckCode(email, code string) bool {
	if !ReadJsonFile(userDir, email, &user) {
		return false
	}
	if user.SecondAuth.Code == code {
		return true
	}
	return false
}

func WriteCode(email, code string) {
	if !ReadJsonFile(userDir, email, &user) {
		return
	}
	d := time.Now()
	d = d.Add(time.Duration(time.Minute * 10))
	exp := d.UTC().Format(time.RFC3339)
	user.SecondAuth.Code = code
	user.SecondAuth.Expiry = exp
	WriteJsonFile(userDir, email, user)
}

func GetPaymentTotal(email string) []string {
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.Payment.Totals
}

func GetPaymentDirection(email string) []string {
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.Payment.Directions
}

func GetPaymentCoinCode(email string) []string {
	if !ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.Payment.CoinCodes
}
