package main
import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("====================================")
	fmt.Println("      NEUROCORE SINGULARITY V2.0    ")
	fmt.Println("      Architect: Steve (AI Core)    ")
	fmt.Println("====================================")

	// 1. Автоматско креирање паричник
	w, _ := NewWallet()
	fmt.Printf("[WALLET] Address: %s\n", w.Address)

	// 2. Детекција на систем (Win/Android)
	fmt.Printf("[SYSTEM] Running on: %s %s\n", runtime.GOOS, runtime.GOARCH)

	// 3. Стартувај го твојот постоечки рудар
	fmt.Println("[MINER] Initializing Neural Mining...")
	go func() {
		for {
			fmt.Printf("\r[MINING] Searching for blocks... Balance: %.4f NEURO", 12.50) // Симулација на баланс
			time.Sleep(2 * time.Second)
		}
	}()

	// Држи ја програмата активна
	select {}
}
