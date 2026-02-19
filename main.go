package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("================================")
	fmt.Println("   NEUROCORE SINGULARITY V2.5   ")
	fmt.Println("================================")
	fmt.Println("1. 🔐 NEW WALLET")
	fmt.Println("2. ⛏️  START MINING")
	fmt.Println("3. ❌ EXIT")
	fmt.Print("\nCHOOSE: ")

	var c int
	fmt.Scan(&c)

	if c == 1 {
		fmt.Printf("\n[!] WALLET CREATED: NC%X\n", rand.Int63())
	} else if c == 2 {
		fmt.Println("\n[+] MINING INITIALIZED...")
		for {
			fmt.Printf("HASH: %d MH/s | BLOCK FOUND!\n", rand.Intn(100)+500)
			time.Sleep(1 * time.Second)
		}
	}
}
