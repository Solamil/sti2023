package main

import (
	"fmt"
	"net/http"
	"strconv"
	"math/rand"
	"time"
	"strings"
	//	"io"
	"text/template"
	"github.com/icza/session"
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
type loginDisplay struct {
	InfoText string
}

type acceptDisplay struct {
	EmailAddress   	string
	InfoText       	string
	PayTotal	string
	PayCoinCode  	string
	PayDirection 	string
}
var indexTemplate *template.Template
var acceptTemplate *template.Template
var loginTemplate *template.Template


func main(){
	//	fmt.Printf("%x", sti2023.Hash(email))
	//fmt.Printf("%x", sti2023.Hash("kukla7@email.cz"))
	generateCode()
	sti2023.CurrencyRates()
	//sti2023.WriteCode(email, "sdfa")
	//mockButton(email)
	//sti2023.CurrencyRates()
	//sti2023.CreatePayment(email, 20.0, "in", "GBP")
	//fmt.Println(sti2023.GetBalances(email))
	//sti2023.AddCurrency(email, "GBP")
	http.HandleFunc("/cover.html", file_handler)
	http.HandleFunc("/index.html", index_handler)
	http.HandleFunc("/accounts", accounts_handler)
	http.HandleFunc("/mock", mock_handler)
	http.HandleFunc("/pay", pay_handler)
	http.HandleFunc("/login", login_handler)
	http.HandleFunc("/logout", logout_handler)
	http.HandleFunc("/", index_handler)

	http.ListenAndServe(":8904", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	// var email string = "michal.kukla@tul.cz"
	
	sess := session.Get(r)
	if sess == nil {
		var i loginDisplay
		i.InfoText = "" 
		loginTemplate, _ = template.ParseFiles("web/login.html")
		loginTemplate.Execute(w, i)
		return
	}
	email := sess.CAttr("Email").(string)
	var info string = ""
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/accept" {
		email := r.FormValue("email")
		totalStr := r.FormValue("total")
		total, _ := strconv.ParseFloat(totalStr, 64)
		direction := r.FormValue("payment_type")
		coinCode := r.FormValue("accounts")
		coinCodePrev := coinCode
		var balance float64 = 0.0
		if sti2023.PreparePayment(email, &balance, &total, &direction, &coinCode) {
			var i acceptDisplay
			i.EmailAddress = email 
			i.PayTotal = fmt.Sprintf("%.2f", total) 
			i.InfoText = ""
			if strings.EqualFold(direction, "IN") {
				i.PayDirection = getHTMLOptionTag(direction, "Příchozí platba")

			} else if strings.EqualFold(direction, "OUT") {
				i.PayDirection = getHTMLOptionTag(direction, "Odchozí platba")
			}
			if strings.EqualFold(coinCode, "CZK") && !strings.EqualFold(coinCodePrev, "CZK") {
				rate := sti2023.GetCurrencyRate(coinCodePrev)
				date := sti2023.GetDate()
				i.InfoText = fmt.Sprintf("%s: %s%s %sKč", date, rate[1], coinCodePrev, rate[0])
			}
			i.PayCoinCode = getHTMLOptionTag(coinCode, coinCode) 
			acceptTemplate, _ = template.ParseFiles("web/accept.html")
			acceptTemplate.Execute(w, i)
			return
		} else {
			info = "Platba se nepodařila provést z důvodu nízkého zůstatku na účtech."
		}

	}
	
	user := sti2023.GetNames(email)
	firstname := user[0]
	lastname := user[1]
	var i indexDisplay
	i.EmailAddress = email
	i.InfoText = info
	i.User = firstname+" "+lastname
	i.AddCurrency = getAddCurrencyHTML(email)
	i.UserCoinCodes = getUserCoinCodesHTML(email)
	i.Accounts = getAccountsHTML(email)
	i.Payments = getPaymentsHTML(email)
	indexTemplate, _ = template.ParseFiles("web/index.html")
	indexTemplate.Execute(w, i)
}

func login_handler(w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if sti2023.CheckPassword(email, password) {
		password = ""
		code := r.PostFormValue("code")

		if (code == "" && sti2023.CheckCode(email, "")) || !sti2023.IsCodeUptodate(email) {
			code := fmt.Sprintf("%d", generateCode())
			infoText, ok := sendCode(email, code)
			if ok {
				sti2023.WriteCode(email, code)
			}

			var i loginDisplay
			i.InfoText = infoText 
			loginTemplate, _ = template.ParseFiles("web/login.html")
			loginTemplate.Execute(w, i)
		} else if code != "" && sti2023.CheckCode(email, code) && sti2023.IsCodeUptodate(email) {
			sess := session.NewSessionOptions(&session.SessOptions{
			    CAttrs: map[string]interface{}{"Email": email},
			})
			session.Add(sess, w)
			sti2023.WriteCode(email, "")
			http.Redirect(w, r, "/", http.StatusFound)
		} else if code != "" && !sti2023.CheckCode(email, code) && !sti2023.IsCodeUptodate(email) {
			var i loginDisplay
			i.InfoText = "⚠️Nepodařilo se přihlásit. Zadaný kód není správný."
			loginTemplate, _ = template.ParseFiles("web/login.html")
			loginTemplate.Execute(w, i)
		}
	} else {
		var i loginDisplay
		i.InfoText = "⚠️Nepodařilo se přihlásit. Heslo nebo emailová adresa není správně."
		loginTemplate, _ = template.ParseFiles("web/login.html")
		loginTemplate.Execute(w, i)
	}
}

func logout_handler(w http.ResponseWriter, r *http.Request) {
	sess := session.Get(r)
	if sess != nil {
		session.Remove(sess, w)		
		sess = nil
	}
	var i loginDisplay	
	i.InfoText = "Jste odhlášen."
	loginTemplate, _ = template.ParseFiles("web/login.html")
	loginTemplate.Execute(w, i)
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web"+r.URL.Path)
}

func pay_handler(w http.ResponseWriter, r *http.Request) {
	total, _ := strconv.ParseFloat(r.FormValue("total"), 64)
	direction := r.FormValue("payment_type")
	coinCode := r.FormValue("accounts")
	email := r.FormValue("email")
	if !sti2023.CreatePayment(email, total, direction, coinCode) {
		var i acceptDisplay
		i.InfoText = "Platba se nezdařila."
		i.PayTotal = fmt.Sprintf("%.2f", total)
		i.PayDirection = direction
		i.PayCoinCode = coinCode
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func mock_handler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	mockButton(email)
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
func sendCode(email, code string) (string, bool) {
	var result string = ""
	var ok bool = false
	if ok = sti2023.Mail(email, code); ok {
		result = fmt.Sprintf("Na emailovou adresu %s Vám byl zaslán ověřovací kód."+
				"Upozornění: Zprává se může nacházet ve složce SPAM.", email)
	} else {
		result = fmt.Sprintf("Na vaši emailovou adresu %s se nepodařilo zaslat ověřovací kód.",
					email)
	}
	return result, ok
}
	
func generateCode() int {
	var max int = 9999
	var min int = 1000
	rand.Seed(time.Now().UnixNano())
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
		if strings.EqualFold(directions[i], "IN") {
			directions[i] = "+"
		}else if strings.EqualFold(directions[i], "out") {
			directions[i] = "-"
		}
		result = fmt.Sprintf("%s<h4 style=\"display: inline-block; margin: 10px;\">%s%s %s</h4>\n", result, directions[i], totals[i], coinCodes[i])

	}
	return result
}

func getHTMLOptionTag(value, symbol string) string {
	var tag string = ""
	tag = fmt.Sprintf("<option value=\"%s\">%s</option>", value, symbol)
	return tag
}
