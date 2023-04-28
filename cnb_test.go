package sti2023

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetCurrencySum(t *testing.T) {
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

	tests := []struct {
		amount   float64
		coinCode string
		exp      string
	}{
		{20.0, "GBP", "537.12"},
		{-1, "GBP", ""},
		{0, "GBPsdaf", ""},
	}

	for _, test := range tests {
		if got := GetCurrencySum(test.amount, test.coinCode); got != test.exp {
			t.Errorf("Expected '%s' but, got '%s'", test.exp, got)
		}
	}
}

func TestGetDate(t *testing.T) {
	cnbDir = "test-cache"
	var exp string = "29.03.2023"
	if got := GetDate(); got != exp {
		t.Errorf("Expected '%s' but, got '%s'", exp, got)
	}
}

func TestGetCoinCodes(t *testing.T) {
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
	tests := []struct {
		exp []string
	}{
		{[]string{"AUD", "BRL", "BGN", "CNY", "DKK", "EUR", "PHP", "HKD", "INR", "IDR", "ISK", "ILS", "JPY", "ZAR", "CAD", "KRW", "HUF", "MYR", "MXN", "XDR", "NOK", "NZD", "PLN", "RON", "SGD", "SEK", "CHF", "THB", "TRY", "USD", "GBP"}},
	}
	for _, test := range tests {
		if got := GetCoinCodes(); len(got) != len(test.exp) {
			t.Errorf("Expected '%s' but, got '%s'", test.exp, got)
		}
	}
}

func TestCleanUpCache(t *testing.T) {
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
