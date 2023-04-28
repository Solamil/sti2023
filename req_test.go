package sti2023

import (
	"testing"
	// "net/http"
	// "net/http/httptest"
)

func TestFailNewrequest(t *testing.T) {
	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprintln(w, `Hello World`)
	// }))

	if got := Newrequest("localhost:8080"); got != "" {

		t.Errorf("Expected '%s' but, got '%s'", "", got)
	}
}
