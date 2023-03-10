package main

import (
	"encoding/json"
	"flag"
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
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port := flagSet.Int("port", 9001, "OAuth server port")
		addr := flagSet.String("addr", "localhost:8080", "Address to proxy")
		err := flagSet.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: parsing flags: %s\n", os.Args[0], err)
			os.Exit(1)
		}

		token, err := waygate.FlowToken("tn.apitman.com", "localhost", *port)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Println("Token: ", token)

		httpClient := &http.Client{}

		privateKey, publicKey, err := waygate.GenerateWireGuardKeys()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

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

		tunRes := &waygate.WireGuardTunnelResponse{}

		err = json.NewDecoder(resp.Body).Decode(tunRes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Println(resp.Status)
		printJson(tunRes)

		wgConfig := tunRes.ClientConfig

		wgConfig.PrivateKey = privateKey

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

		caddyConfig := waygate.NewCaddyHttpConfig()
		caddyRoute := waygate.CaddyRoute{
			Match: []interface{}{
				struct {
					Host []string `json:"host"`
				}{
					Host: tunRes.Domains,
				},
			},
			Handle: []waygate.CaddyHandler{
				waygate.CaddyHandler{
					Handler: "reverse_proxy",
					Upstreams: []waygate.CaddyUpstream{
						waygate.CaddyUpstream{
							Dial: *addr,
						},
					},
				},
			},
		}

		caddyConfig.Apps.Http.Servers.Waygate.Listen = []string{":5757"}
		caddyConfig.Apps.Http.Servers.Waygate.Routes = []waygate.CaddyRoute{
			caddyRoute,
		}

		printJson(caddyConfig)

		caddyConfigBytes, err := json.MarshalIndent(caddyConfig, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = os.WriteFile("./caddy.json", caddyConfigBytes, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		cmd = exec.Command("caddy", "run", "--config", "./caddy.json")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, string(output))
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func printJson(data interface{}) {
	d, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(d))
}
