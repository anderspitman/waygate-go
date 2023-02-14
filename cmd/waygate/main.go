package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

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

		privateKey := "wOEXYf4xNqtrFokL+LATj8EYQK+1ughMDqXnvlbj72Y="
		publicKey := "fPy5iEIhAQxIlurDiY4W+qEvXsF/t1a/koapEkVbrDc="

		pubKeyEscaped := url.QueryEscape(publicKey)
		url := fmt.Sprintf("https://tn.apitman.com/waygate/open?type=wireguard&token=%s&public-key=%s", token, pubKeyEscaped)
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

		wgConfig := &waygate.WireGuardConfig{}

		err = json.NewDecoder(resp.Body).Decode(wgConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		wgConfig.PrivateKey = privateKey

		fmt.Println(resp.Status)

		fmt.Println(wgConfig)

		err = os.WriteFile("/etc/wireguard/waygate0.conf", []byte(wgConfig.String()), 0600)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		cmd := exec.Command("wg-quick", "up", "waygate0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, string(output))
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}
