package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UrlToCheck struct {
	Url     string
	IsValid bool
}

func (u *UrlToCheck) check() {
	nReq, _ := http.NewRequest("GET", u.Url, nil)

	if resp, curlErr := (&http.Client{}).Do(nReq); curlErr == nil && resp.StatusCode == http.StatusOK {
		u.IsValid = true
	}
}

func main() {
	http.HandleFunc("/", handleInputData)

	err := http.ListenAndServe(":666", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handleInputData(w http.ResponseWriter, req *http.Request) {
	cType := req.Header.Get("Content-Type")

	if cType != "application/json" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, "Can't run io.ReadAll", http.StatusInternalServerError)
	}

	urlList := []string{}
	response := []interface{}{}
	json.Unmarshal(body, &urlList)

	for _, url := range urlList {
		urlData := UrlToCheck{url, false}
		urlData.check()

		response = append(response, urlData)
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
