package main
import (
	"fmt"
	"math/rand"
	"time"
)
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== NEUROCORE V2.5 ACTIVE ===")
	addr := "NC" + fmt.Sprintf("%X", rand.Int63())
	fmt.Println("Your New Wallet:", addr)
	fmt.Println("Mining starting now...")
	for {
		fmt.Printf("Mined 0.005 NC | Hash: %d MH/s\n", rand.Intn(100)+500)
		time.Sleep(2 * time.Second)
	}
}
