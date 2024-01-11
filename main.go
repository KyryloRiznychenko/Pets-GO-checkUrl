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

func (u *UrlToCheck) check(c chan<- UrlToCheck) {
	nReq, _ := http.NewRequest("GET", u.Url, nil)

	if resp, curlErr := (&http.Client{}).Do(nReq); curlErr == nil && resp.StatusCode == http.StatusOK {
		u.IsValid = true
	}

	c <- *u
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
	json.Unmarshal(body, &urlList)
	ch := make(chan UrlToCheck, len(urlList))
	for _, url := range urlList {
		go (&UrlToCheck{Url: url}).check(ch)
	}

	response := make([]UrlToCheck, len(urlList))
	for i := 0; i < len(urlList); i++ {
		response = append(response, <-ch)
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
