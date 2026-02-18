package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	mempoolCh    chan Transaction
	newBlockCh   chan Block
	minerStopCh  chan struct{}
	minerRunning int32
)

func StartMiner() bool {
	if !atomic.CompareAndSwapInt32(&minerRunning, 0, 1) {
		return false
	}
	if mempoolCh == nil {
		mempoolCh = make(chan Transaction, 100)
	}
	if newBlockCh == nil {
		newBlockCh = make(chan Block)
	}
	minerStopCh = make(chan struct{})
	go Miner(mempoolCh, newBlockCh, minerStopCh)
	return true
}

func StopMiner() bool {
	if !atomic.CompareAndSwapInt32(&minerRunning, 1, 0) {
		return false
	}
	if minerStopCh != nil {
		close(minerStopCh)
		minerStopCh = nil
	}
	return true
}

func walletNewHandler(wr http.ResponseWriter, r *http.Request) {
	wlt, err := NewWallet()
	if err != nil {
		http.Error(wr, err.Error(), 500)
		return
	}
	out := map[string]string{"address": wlt.Address, "private_key": wlt.PrivateKeyHex()}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(out)
}

func txHandler(wr http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(wr, "method", 405)
		return
	}
	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(wr, err.Error(), 400)
		return
	}
	if mempoolCh == nil {
		mempoolCh = make(chan Transaction, 100)
	}
	select {
	case mempoolCh <- tx:
	default:
	}
	wr.WriteHeader(200)
}

func mineStartHandler(wr http.ResponseWriter, r *http.Request) {
	started := StartMiner()
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]bool{"started": started})
}

func mineStopHandler(wr http.ResponseWriter, r *http.Request) {
	stopped := StopMiner()
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]bool{"stopped": stopped})
}

func statusHandler(wr http.ResponseWriter, r *http.Request) {
	height := len(Blockchain)
	total := atomic.LoadUint64(&globalHashCounter)
	status := map[string]interface{}{
		"height":        height,
		"hashes_total":  total,
		"network":       NetworkStatus(),
		"miner_running": atomic.LoadInt32(&minerRunning) == 1,
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(status)
}

func marketListHandler(wr http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			Seller, Asset string
			Price         float64
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(wr, err.Error(), 400)
			return
		}
		l := Market.CreateListing(req.Seller, req.Asset, req.Price)
		wr.Header().Set("Content-Type", "application/json")
		json.NewEncoder(wr).Encode(l)
		return
	}
	// GET
	list := Market.ListListings()
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(list)
}

func marketBuyHandler(wr http.ResponseWriter, r *http.Request) {
	var req struct{ ListingID, Buyer string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(wr, err.Error(), 400)
		return
	}
	l, err := Market.Buy(req.ListingID, req.Buyer)
	if err != nil {
		http.Error(wr, err.Error(), 400)
		return
	}
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(l)
}

func main() {
	// handle CLI flags (inspect/submit)
	if HandleCLI() {
		return
	}

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

	// initialize channels
	mempoolCh = make(chan Transaction, 100)
	newBlockCh = make(chan Block)

	// start P2P server
	go StartP2PServer("40404")

	// block handler: listen for mined blocks
	go func() {
		for b := range newBlockCh {
			fmt.Println("Mined block", b.Index, "hash", b.Hash[:12], "nonce", b.Nonce)
			AddBlock(b)
			data, _ := json.Marshal(b)
			Broadcast(Message{Type: "NewBlock", Data: data})
		}
	}()

	// simple dashboard printer
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

	// provide a sample transaction every 30s to keep mempool active
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			// use a throwaway wallet
			w, _ := NewWallet()
			tx := Transaction{From: w.Address, To: w.Address, Amount: 0.001}
			select {
			case mempoolCh <- tx:
			default:
			}
			data, _ := json.Marshal(tx)
			Broadcast(Message{Type: "NewTransaction", Data: data})
		}
	}()

	// HTTP API routes
	http.HandleFunc("/wallet/new", walletNewHandler)
	http.HandleFunc("/tx", txHandler)
	http.HandleFunc("/mine/start", mineStartHandler)
	http.HandleFunc("/mine/stop", mineStopHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/market/list", marketListHandler)
	http.HandleFunc("/market/listings", marketListHandler)
	http.HandleFunc("/market/buy", marketBuyHandler)

	go func() {
		fmt.Println("HTTP API listening on :8080")
		http.ListenAndServe(":8080", nil)
	}()

	// start miner automatically
	StartMiner()

	// wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	StopMiner()
	fmt.Println("Shutting down...")
}
