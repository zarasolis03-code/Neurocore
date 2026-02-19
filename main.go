package main
import (
	"fmt"
	"math/rand"
	"time"
)
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== NEUROCORE V2.1 ===")
	fmt.Println("1. Create Wallet\n2. Start Mining")
	var c int
	fmt.Scan(&c)
	if c == 1 {
		fmt.Printf("Wallet: NC%X\n", rand.Int63())
	} else {
		fmt.Println("Mining started...")
		for {
			fmt.Printf("Hash: %d MH/s\n", rand.Intn(100)+400)
			time.Sleep(2 * time.Second)
		}
	}
}
