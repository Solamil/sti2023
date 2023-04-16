package main

import (
	"fmt"
	"net/http"
	"strconv"
	"math/rand"
	"time"
	//	"io"
	"github.com/Solamil/sti2023"
)

func main(){
	fmt.Printf("%x", sti2023.Hash("michal.kukla@tul.cz"))
	generateCode()

	sti2023.WriteCode("michal.kukla@tul.cz", "sdfa")
	mockButton("michal.kukla@tul.cz")
	//sti2023.CurrencyRates()
	//sti2023.CreatePayment("michal.kukla@tul.cz", 20.0, "in", "GBP")
	//fmt.Println(sti2023.GetBalances("michal.kukla@tul.cz"))
	//sti2023.AddCurrency("michal.kukla@tul.cz", "GBP")
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/index.html", index_handler)
	http.ListenAndServe(":8904", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {


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
