package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36"

var myClient = &http.Client{Timeout: 10 * time.Second}

// https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
func getJson(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", userAgent)

	r, err := myClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	return json.NewDecoder(r.Body).Decode(target)
}

func postJson(url string, target interface{}, body []byte) error {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("Content-type", "application/json")
	r, err := myClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	return json.NewDecoder(r.Body).Decode(target)
}