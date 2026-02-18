package main
import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)
func main() {
    fmt.Println("\033[35mNEUROCORE SINGULARITY V2.0\033[0m")
    fmt.Println("Created by Steve & Neurocore AI\n")
    LoadChain()
    w, _ := NewWallet()
    blockCh := make(chan Block)
    stopCh := make(chan struct{})
    go Miner(blockCh, stopCh)
    go StartP2P()
    go func() {
        for b := range blockCh {
            AddBlock(b)
            fmt.Printf("\r\033[K\033[32m[Mined]\033[0m Block #%d | Balance: %.2f NEURO", b.Index, GetBalance(w.Address))
        }
    }()
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        for range ticker.C {
            fmt.Printf("\n\033[36m[AI Chat]:\033[0m %s\n", GetAIChatResponse("status", GetBalance(w.Address)))
        }
    }()
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs
    close(stopCh)
    fmt.Println("\nShutdown complete.")
}
