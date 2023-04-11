package sti2023

import (
	"fmt"
	"strconv"
)

type User struct {
	Password  string   `json:"password"`
	Balances  []string `json:"balances"`
	CoinCodes []string `json:"coinCodes"`
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Payment  struct {
		Directions []string `json:"directions"`
		Totals     []string `json:"totals"`
		CoinCodes  []string `json:"coinCodes"`
	} `json:"Payments"`
	VerifyCode string `json:"verify"`
}

var userDir string = "users"

func CreatePayment(email string, total float64, direction string, coinCode string) bool {
	
	index := getIndex(GetUserCoinCodes(email), coinCode)
	if index < 0 || total <= 0 {
		return false
	}

	balances := GetBalances(email)
	balance, _ := strconv.ParseFloat(balances[index], 64)
		
	if direction == "in" {
		balance += total
	} else if direction == "out" && balance < total {

		valueCzk := GetCurrencySum(total, coinCode)
		total, _ = strconv.ParseFloat(valueCzk, 64)
		index = 0 // CZK
		coinCode = "CZK"
		balance, _ = strconv.ParseFloat(balances[index], 64)
		fmt.Println(balance)
		if balance < total {
			return false
		}
		balance -= total

	} else if direction == "out" && balance >= total {
		balance -= total	
	} else {
		return false	
	}
	return AddPayment(email, balance, total, direction, coinCode) 
}

func AddCurrency(email string, code string) {
	var user User
	if ! ReadJsonFile(userDir, email, &user) {
		return
	}

	if ! IsExistCode(code) {
		return
	}
	user.Balances = append(user.Balances, "0.0")
	user.CoinCodes = append(user.CoinCodes, code) 	
	WriteJsonFile(userDir, email, user)
}

func GetBalance(email string, code string) string {
	var result string = ""
	var index int = getIndex(GetUserCoinCodes(email), code)
	if index < 0 {
		return result	
	}
	result = GetBalances(email)[index]
	return result 
}

func GetBalances(email string) []string {
	var user User
	if ! ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.Balances
}

func GetUserCoinCodes(email string) []string {
	var user User
	if ! ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return user.CoinCodes
}

func GetNames(email string) []string {
	var user User
	if ! ReadJsonFile(userDir, email, &user) {
		return []string{}
	}
	return []string{user.FirstName, user.LastName}
}

func IsCorrectCode(email, code string) bool {
	var user User
	ReadJsonFile(userDir, email, &user)
	if user.VerifyCode == code {
		WriteCode(email, "")
		user.VerifyCode = ""
		
		return WriteJsonFile(userDir, email, user)
	}
	return false
}

func AddPayment(email string, balance, total float64, direction, coinCode string) bool {
	var user User
	if ! ReadJsonFile(userDir, email, &user) {
		return false
	}
	index := getIndex(GetUserCoinCodes(email), coinCode)
	strTotal := fmt.Sprintf("%.2f", total)
	strBalance := fmt.Sprintf("%.2f", balance)
	user.Balances[index] = strBalance
	user.Payment.Directions = append(user.Payment.Directions, direction)
	user.Payment.Totals = append(user.Payment.Totals, strTotal)
	user.Payment.CoinCodes = append(user.Payment.CoinCodes, coinCode)

	return WriteJsonFile(userDir, email, user)
}

func WriteCode(email, code string) {
	var user User
	ReadJsonFile(userDir, email, &user)
	user.VerifyCode = code
	WriteJsonFile(userDir, email, user)
}

func GetPaymentsHTML(email string) []string {
	return []string{}
}
