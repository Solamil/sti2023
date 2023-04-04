package main

import (
	//	"fmt"
	"net/http"
	//	"io"
	//	"github.com/Solamil/sti2023"
)

func main(){
	//	sti2023.CreatePayment("michal.kukla@tul.cz", 0.0, "f", "fdsf")
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/index.html", index_handler)
	http.ListenAndServe(":8904", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {


}


