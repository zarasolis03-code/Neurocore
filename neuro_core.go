package main
import (
	"fmt"
	"math/rand"
	"time"
)
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(">> NEUROCORE SINGULARITY V2.1 <<")
	fmt.Printf("NEW WALLET GENERATED: NC%X\n", rand.Int63())
	fmt.Println("System initializing...")
	time.Sleep(2 * time.Second)
	for {
		fmt.Printf("MINING ACTIVE | Hashrate: %d MH/s | Block Found!\n", rand.Intn(100)+450)
		time.Sleep(1 * time.Second)
	}
}
