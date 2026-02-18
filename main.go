package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	// handle CLI flags (inspect/submit)
	if HandleCLI() {
		return
	}

	// initialize wallet
	w, err := NewWallet()
	if err != nil {
		panic(err)
	}
	fmt.Println("Wallet Address:", w.Address)
	fmt.Println("Private Key (hex):", w.PrivateKeyHex())

	// create genesis if needed
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		createGenesisBlock()
		fmt.Println("Created genesis block and wrote to", historyFile)
	} else {
		// try to load existing history
		data, err := os.ReadFile(historyFile)
		if err == nil {
			json.Unmarshal(data, &Blockchain)
		}
	}

	// channels
	mempoolCh := make(chan Transaction, 100)
	newBlockCh := make(chan Block)
	stopCh := make(chan struct{})

	// start miner
	go Miner(mempoolCh, newBlockCh, stopCh)

	// start P2P server
	go StartP2PServer("40404")

	// handle mined blocks
	go func() {
		for b := range newBlockCh {
			fmt.Println("Mined block", b.Index, "hash", b.Hash[:12], "nonce", b.Nonce)
			AddBlock(b)
			// broadcast
			data, _ := json.Marshal(b)
			Broadcast(Message{Type: "NewBlock", Data: data})
		}
	}()

	// simple dashboard
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		var lastCount uint64
		for range ticker.C {
			height := len(Blockchain)
			total := atomic.LoadUint64(&globalHashCounter)
			hs := total - lastCount
			lastCount = total
			fmt.Printf("[Dashboard] Height=%d  H/s=%d  %s\n", height, hs/3, NetworkStatus())
		}
	}()

	// example: push a dummy tx every 10 seconds
	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()
		for range tick.C {
			tx := Transaction{From: w.Address, To: w.Address, Amount: 0.01}
			mempoolCh <- tx
			data, _ := json.Marshal(tx)
			Broadcast(Message{Type: "NewTransaction", Data: data})
		}
	}()

	// wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stopCh)
	fmt.Println("Shutting down...")
}
