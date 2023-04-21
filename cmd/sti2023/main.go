package main

import (
	"fmt"
	"net/http"
	"strconv"
	"net/mail"
	"math/rand"
	"time"
	//	"io"
	"text/template"
	"github.com/Solamil/sti2023"
)

type indexDisplay struct {
	EmailAddress   string
	InfoText       string
	User	       string
	AddCurrency    string
	UserCoinCodes  string
	Accounts       string
	Payments       string
}
var indexTemplate *template.Template



func main(){
	//	fmt.Printf("%x", sti2023.Hash(email))
	generateCode()

	//sti2023.WriteCode(email, "sdfa")
	//mockButton(email)
	//sti2023.CurrencyRates()
	//sti2023.CreatePayment(email, 20.0, "in", "GBP")
	//fmt.Println(sti2023.GetBalances(email))
	//sti2023.AddCurrency(email, "GBP")
	http.HandleFunc("/index.html", index_handler)
	http.HandleFunc("/accounts", accounts_handler)
	http.HandleFunc("/mock", mock_handler)
	http.HandleFunc("/pay", pay_handler)
	http.HandleFunc("/", index_handler)

	http.ListenAndServe(":8904", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	var email string = "michal.kukla@tul.cz"

	user := sti2023.GetNames(email)
	firstname := user[0]
	lastname := user[1]
	var i indexDisplay
	i.EmailAddress = email
	i.InfoText = ""
	i.User = firstname+" "+lastname
	i.AddCurrency = getAddCurrencyHTML(email)
	i.UserCoinCodes = getUserCoinCodesHTML(email)
	i.Accounts = getAccountsHTML(email)
	i.Payments = getPaymentsHTML(email)
	indexTemplate, _ = template.ParseFiles("web/index.html")
	indexTemplate.Execute(w, i)
}

func pay_handler(w http.ResponseWriter, r *http.Request) {
	//total := r.FormValue("total")
	//direction := r.FormValue("payment_type")
	//coin := r.FormValue("accounts")
	//email := r.FormValue("email")
	//sti2023.CreatePayment(email, total, direction, coin)
}

func mock_handler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	mockButton(email)
	//fmt.Println(r.FormValue("email"))
	http.Redirect(w, r, "/", http.StatusFound)
}

func accounts_handler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	coin := r.FormValue("currencies")
	sti2023.AddCurrency(email, coin)	
	http.Redirect(w, r, "/", http.StatusFound)
}

func mockButton(email string) bool {
	rand.Seed(time.Now().UnixNano())
	balances := sti2023.GetBalances(email)
	coinIndex := rand.Intn(len(balances))
	coinCode := sti2023.GetUserCoinCodes(email)[coinIndex]
	balance, _ := strconv.ParseFloat(balances[coinIndex], 64)
	
	rand.Seed(time.Now().UnixNano())

	direction := "out"
	if balance < 25 {
		direction = "in"
	}

	var total float64 = 0.0 
	if balance > 1 {
		total += float64(rand.Intn(int(balance)))
	}
	total += rand.Float64() 
	return sti2023.CreatePayment(email, total, direction, coinCode) 
}
	
func generateCode() int {
	var max int = 9999
	var min int = 1000
	var number int = rand.Intn(max-min) + min
	return number
}

func getAddCurrencyHTML(email string) string {
	var result string = ""

	userCoinCodes := sti2023.GetUserCoinCodes(email)
	coinCodes := sti2023.GetCoinCodes()
	
	for _, v := range coinCodes {
		if i := sti2023.GetIndex(userCoinCodes, v); i < 0 {
			result = fmt.Sprintf("%s%s\n", result, getHTMLOptionTag(v, v))
		}
	}
	return result
}

func getUserCoinCodesHTML(email string) string {
	var result string = ""

	for _, v := range sti2023.GetUserCoinCodes(email) {
		result = fmt.Sprintf("%s%s\n", result, getHTMLOptionTag(v,v))
	}
	return result
}

func getAccountsHTML(email string) string {
	var result string = ""

	userCoinCodes := sti2023.GetUserCoinCodes(email)

	for i, balance := range sti2023.GetBalances(email) {
		if len(userCoinCodes) <= i {
			break
		}
		result = fmt.Sprintf("%s<p style=\"display: inline-block\">%s %s</p>", result,
					balance, userCoinCodes[i])
	}
	return result
}

func getPaymentsHTML(email string) string {
	var result string = ""
	totals := sti2023.GetPaymentTotal(email)		
	directions := sti2023.GetPaymentDirection(email)
	coinCodes := sti2023.GetPaymentCoinCode(email)
	if len(totals) != len(directions) || len(totals) != len(coinCodes) {
		return result
	}
	for i := len(totals)-1; i >= 0; i-- {
		if directions[i] == "in" {
			directions[i] = "+"
		}else if directions[i] == "out" {
			directions[i] = "-"
		}
		result = fmt.Sprintf("%s<h4 style=\"display: inline-block; margin: 10px;\">%s%s %s</h4>\n", result, directions[i], totals[i], coinCodes[i])

	}
	return result
}

func validMailAddress(address string) (string, bool) {
    addr, err := mail.ParseAddress(address)
    if err != nil {
        return "", false
    }
    return addr.Address, true
}

func getHTMLOptionTag(value, symbol string) string {
	var tag string = ""
	tag = fmt.Sprintf("<option value=\"%s\">%s</option>", value, symbol)
	return tag
}
