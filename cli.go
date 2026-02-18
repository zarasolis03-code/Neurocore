package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// HandleCLI checks for -submit or -inspect and runs the action, returning true if the program should exit.
func HandleCLI() bool {
	submit := flag.Bool("submit", false, "Submit a sample transaction from the generated wallet")
	inspect := flag.Bool("inspect", false, "Print the current chain to stdout")
	flag.Parse()

	if *submit {
		w, err := NewWallet()
		if err != nil {
			fmt.Println("wallet error:", err)
			os.Exit(1)
		}
		tx := Transaction{From: w.Address, To: w.Address, Amount: 0.1}
		data, _ := json.Marshal(tx)
		fmt.Println(string(data))
		return true
	}
	if *inspect {
		data, _ := json.MarshalIndent(Blockchain, "", "  ")
		fmt.Println(string(data))
		return true
	}
	return false
}
