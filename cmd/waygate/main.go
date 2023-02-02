package main

import (
	"fmt"
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
	}
}
