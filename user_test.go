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

		{email, 140.0, "in", "CZK", true},
		{email, 1.25, "in", "GBP", true},
		{email, 2.0, "out", "GBP", true},
		{email, 1.25, "out", "GBP", true},
		{email, 86.29, "out", "CZK", true},

		{email, 10.0, "out", "GBP", false},
		{wrongEmail, 20.0, "in", "CZK", false},
		{email, 20.0, "out", "ABC", false},
		{email, 0.0, "in", "CZK", false},
	}

	for _, test := range tests {
		if got := CreatePayment(test.email, test.total, test.direction, test.coinCode); test.exp != got {
			t.Errorf("Expected '%t' but, got '%t',\n %.2f %s %s", test.exp,
				got, test.total, test.direction, test.coinCode)
		}
	}

	var user User
	ReadJsonFile(userDir, email, &user)
	var exp string = "0.00"
	for _, balance := range user.Balances {
		if balance != exp {
			t.Errorf("Expected '%s' but, got '%s'", exp, balance)

		}
	}
}

//  func TestCleanUpCache1(t *testing.T) {
//  	var cache_dir string = "test-cache"
//  	dirRead, _ := os.Open(cache_dir)
//  	dirFiles, _ := dirRead.Readdir(0)
//  	for index := range dirFiles {
//  		file := dirFiles[index]
//  		filename := file.Name()
//  		if err := os.Remove(cache_dir + "/" + filename); err != nil {
//  			t.Errorf("error %s", err)
//  		}
//  	}
//  	if err := os.Remove(cache_dir + "/"); err != nil {
//  		t.Errorf("error %s", err)
//  	}
//  }

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
