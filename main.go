package main
import (
	"fmt"
	"math/rand"
	"time"
)
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("================================")
	fmt.Println("   NEUROCORE SINGULARITY V2.1   ")
	fmt.Println("================================")
	fmt.Printf("WALLET: NC%X\n", rand.Int63())
	fmt.Println("MINING STATUS: ACTIVE")
	for {
		fmt.Printf("Mined 0.0001 NC | Hash: %d MH/s\n", rand.Intn(100)+450)
		time.Sleep(2 * time.Second)
	}
}
