package sti2023

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreatePayment(t *testing.T) {
	userDir = "test-cache"
	cnbDir = "test-cache"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open("test-data/devizovy_trh.json")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		byteRates, _ := io.ReadAll(file)
		w.Write(byteRates)
	}))
	cnbUrl = ts.URL

	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)
	tests := []struct {
		email     string
		total     float64
		direction string
		coinCode  string
		exp       bool
	}{

		{email, 140.0, "IN", "CZK", true},
		{email, 1.25, "IN", "GBP", true},
		{email, 2.0, "OUT", "GBP", true},
		{email, 1.25, "OUT", "GBP", true},
		{email, 86.29, "OUT", "CZK", true},

		{email, 10.0, "OUT", "GBP", false},
		{wrongEmail, 20.0, "IN", "CZK", false},
		{email, 20.0, "OUT", "ABC", false},
		{email, 0.0, "IN", "CZK", false},
	}

	for _, test := range tests {
		if got := CreatePayment(test.email, test.total, test.direction, test.coinCode); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %.2f %s %s", test.exp,
				got, test.total, test.direction, test.coinCode)
		}
	}

	if len(GetPaymentTotal(email)) != 5 || len(GetPaymentDirection(email)) != 5 ||
		len(GetPaymentCoinCode(email)) != 5 {
		t.Errorf(`Expected 5 length for address %s, but got total: %d, direction: %d,
		coinCode: %d`, wrongEmail, len(GetPaymentTotal(wrongEmail)),
			len(GetPaymentDirection(wrongEmail)), len(GetPaymentCoinCode(wrongEmail)))
	}
	if len(GetPaymentTotal(wrongEmail)) != 0 || len(GetPaymentDirection(wrongEmail)) != 0 ||
		len(GetPaymentCoinCode(wrongEmail)) != 0 {
		t.Errorf(`Expected 0 length for address %s, but got total: %d, direction: %d,
		coinCode: %d`, wrongEmail, len(GetPaymentTotal(wrongEmail)),
			len(GetPaymentDirection(wrongEmail)), len(GetPaymentCoinCode(wrongEmail)))
	}
	var exp string = "0.00"
	for _, balance := range GetBalances(email) {
		if balance != exp {
			t.Errorf("Expected '%s' but, got '%s'", exp, balance)

		}
	}
}

func TestPreparePayment(t *testing.T) {
	userDir = "test-cache"
	cnbDir = "test-cache"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open("test-data/devizovy_trh.json")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		byteRates, _ := io.ReadAll(file)
		w.Write(byteRates)
	}))
	cnbUrl = ts.URL

	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	var user User
	user.Password = "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578"
	user.Balances = append(user.Balances, "140.0")
	user.CoinCodes = append(user.CoinCodes, "CZK")
	user.Balances = append(user.Balances, "1.25")
	user.CoinCodes = append(user.CoinCodes, "GBP")
	user.FirstName = "Michal"
	user.LastName = "Kukla"

	WriteJsonFile(userDir, email, user)
	tests := []struct {
		email        string
		balance      float64
		total        float64
		direction    string
		coinCode     string
		exp          bool
		expBalance   float64
		expTotal     float64
		expDirection string
		expCoinCode  string
	}{

		{email, 0.0, 140.0, "IN", "CZK", true, 280.0, 140.0, "IN", "CZK"},
		{email, 0.0, 1.25, "IN", "GBP", true, 2.5, 1.25, "IN", "GBP"},
		{email, 0.0, 2.0, "OUT", "GBP", true, 86.29, 53.71, "OUT", "CZK"},
		{email, 0.0, 1.25, "OUT", "GBP", true, 0.0, 1.25, "OUT", "GBP"},
		{email, 0.0, 86.29, "OUT", "CZK", true, 53.71, 86.29, "OUT", "CZK"},

		{email, 0.0, 10.0, "OUT", "GBP", false, 140.0, 268.56, "OUT", "CZK"},
		{wrongEmail, 0.0, 20.0, "IN", "CZK", false, 0.0, 20.0, "IN", "CZK"},
		{email, 0.0, 20.0, "OUT", "ABC", false, 0.0, 20.0, "OUT", "ABC"},
		{email, 0.0, 0.0, "IN", "CZK", false, 0.0, 0.0, "IN", "CZK"},
		{email, 0.0, 400.0, "OUT", "CZK", false, 0.0, 400.0, "OUT", "CZK"},
	}

	for _, test := range tests {
		if got := PreparePayment(test.email, &test.balance, &test.total, &test.direction, &test.coinCode); test.exp != got ||
			test.total != test.expTotal || test.direction != test.expDirection || test.coinCode != test.expCoinCode {
			t.Errorf(`Expected '%t' but, got '%t',
			Got: %.2f, %.2f, %s, %s,
			Expected: %.2f, %.2f, %s, %s`, test.exp, got,
				test.balance, test.total, test.direction, test.coinCode,
				test.expBalance, test.expTotal, test.expDirection, test.expCoinCode)
		}
	}
}

func TestAddPayment(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"

	setDefaultUser(userDir, email)
	tests := []struct {
		email     string
		balance   float64
		total     float64
		direction string
		coinCode  string
		exp       bool
	}{

		{email, 0.0, 140.0, "IN", "CZK", true},
		{wrongEmail, 0.0, 140.0, "IN", "CZK", false},
		{email, 0.0, 20.0, "OUT", "ABC", false},
	}
	for _, test := range tests {
		if got := addPayment(test.email, test.balance, test.total, test.direction, test.coinCode); got != test.exp {
			t.Errorf("Expected '%t' but, got '%t'", test.exp, got)
		}
	}
}

func TestAddCurrency(t *testing.T) {
	userDir = "test-cache"
	cnbDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)

	tests := []struct {
		email    string
		coinCode string
		exp      bool
	}{

		{email, "ABC", false},
		{wrongEmail, "CZK", false},
		{wrongEmail, "EUR", false},
		{email, "EUR", true},
		{email, "EUR", false},
	}

	for _, test := range tests {
		if got := AddCurrency(test.email, test.coinCode); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %s %s", test.exp,
				got, test.email, test.coinCode)
		}
	}

	var exp string = "EUR"
	for _, coinCode := range GetUserCoinCodes(email) {
		if coinCode == exp {
			return
		}
	}
	t.Errorf("Expected '%s' to find in array, but cannot find it", exp)

}

func TestCheckPassword(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)
	tests := []struct {
		email string
		value string
		exp   bool
	}{

		{wrongEmail, "testtest", false},
		{email, "tedsfa", false},
		{email, "testest", false},
		{email, "testtest", true},
	}
	for _, test := range tests {
		if got := CheckPassword(test.email, test.value); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %s %s", test.exp,
				got, test.email, test.value)
		}
	}
}

func TestVerifyCode(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)
	tests := []struct {
		email string
		code  string
		exp   bool
	}{

		{wrongEmail, "testtest", false},
		{email, "", true},
		{email, "testtest", true},
	}
	// Test empty string for the start
	if got := CheckCode(email, ""); !got {
		t.Errorf("Expected '%t' but, got '%t',\n %s %s", true,
			got, email, "")
	}

	for _, test := range tests {
		WriteCode(test.email, test.code)
		if got := CheckCode(test.email, test.code); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %s %s", test.exp,
				got, test.email, test.code)
		}
	}
	// Test False
	if got := CheckCode(email, "fff"); got {
		t.Errorf("Expected '%t' but, got '%t',\n %s %s", false,
			got, email, "")
	}
}

func TestIsCodeUptodate(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)
	tests := []struct {
		email string
		code  string
		exp   bool
	}{
		{wrongEmail, "testtest", false},
		{email, "", true},
		{email, "testtest", true},
	}

	// Test empty string for the start
	if got := IsCodeUptodate(email); got {
		t.Errorf("Expected '%t' but, got '%t',\n %s %s", false,
			got, email, "")
	}
	for _, test := range tests {
		WriteCode(test.email, test.code)
		if got := IsCodeUptodate(test.email); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %s %s", test.exp,
				got, test.email, test.code)
		}
	}
}

func TestGetNames(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)

	tests := []struct {
		email string
		exp   []string
	}{
		{wrongEmail, []string{}},
		{email, []string{"Michal", "Kukla"}},
	}

	for _, test := range tests {
		if got := GetNames(test.email); len(test.exp) != len(got) {
			t.Errorf("Expected '%v' but, got '%v',\n %s", test.exp,
				got, test.email)
		}
	}
}

func TestGetBalance(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)

	tests := []struct {
		email string
		code  string
		exp   string
	}{
		{wrongEmail, "CZK", ""},
		{email, "CZK", "0.0"},
		{email, "ABC", ""},
	}

	for _, test := range tests {
		if got := GetBalance(test.email, test.code); test.exp != got {
			t.Errorf("Expected '%s' but, got '%s',\n %s %s", test.exp,
				got, test.email, test.code)
		}
	}
}

func TestGetBalances(t *testing.T) {
	userDir = "test-cache"
	var email string = "michal.kukla@tul.cz"
	var wrongEmail string = "sdkfafa@skdafj.cz"
	setDefaultUser(userDir, email)

	tests := []struct {
		email string
		exp   []string
	}{
		{wrongEmail, []string{}},
		{email, []string{"0.0", "0.0"}},
	}

	for _, test := range tests {
		if got := GetBalances(test.email); len(test.exp) != len(got) {
			t.Errorf("Expected '%v' but, got '%v',\n %s", test.exp,
				got, test.email)
		}
	}
}

func TestCleanUpCache1(t *testing.T) {
	var cache_dir string = "test-cache"
	dirRead, _ := os.Open(cache_dir)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		file := dirFiles[index]
		filename := file.Name()
		if err := os.Remove(cache_dir + "/" + filename); err != nil {
			t.Errorf("error %s", err)
		}
	}
	if err := os.Remove(cache_dir + "/"); err != nil {
		t.Errorf("error %s", err)
	}
}

func setDefaultUser(dir, email string) {
	var user User
	user.Password = "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578"
	user.Balances = append(user.Balances, "0.0")
	user.CoinCodes = append(user.CoinCodes, "CZK")
	user.Balances = append(user.Balances, "0.0")
	user.CoinCodes = append(user.CoinCodes, "GBP")
	user.FirstName = "Michal"
	user.LastName = "Kukla"

	WriteJsonFile(dir, "michal.kukla@tul.cz", user)
}
