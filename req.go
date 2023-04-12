package sti2023

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func Newrequest(url string) string {
	var answer string = ""
	//	t := time.Now().Add(2 * time.Second)
	//	ctx, cancel := context.WithCancel(context.TODO())
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 2 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	reqm, _ := http.NewRequest("GET", url, nil)
	reqm.Header.Set("Content-Type", "text/html")
	content, err := client.Do(reqm)
	if err != nil {
		fmt.Println(err)
		if content != nil {
			fmt.Println("statusCode: ", content.StatusCode)
		}
		return answer
	} else if content.StatusCode >= 400 {
		return answer
	}

	value, err := io.ReadAll(content.Body)
	if err != nil {
		fmt.Println(err)
		return answer
	}
	answer = string(value)
	return answer
}
