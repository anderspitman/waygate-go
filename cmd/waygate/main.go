package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anderspitman/waygate-go"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, os.Args[0]+": Need a command")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "client":
		token, err := waygate.FlowToken("tn.apitman.com", "localhost:9001")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Println("Token: ", token)

		httpClient := &http.Client{}

		publicKey := "fPy5iEIhAQxIlurDiY4W+qEvXsF/t1a/koapEkVbrDc="

		url := fmt.Sprintf("https://tn.apitman.com/waygate/open?type=wireguard&token=%s&public-key=%s", token, publicKey)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer resp.Body.Close()

		fmt.Println(resp.Status)
	}
}
