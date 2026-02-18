package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Signature []byte  `json:"signature"`
}

type Block struct {
	Index         int           `json:"index"`
	Timestamp     int64         `json:"timestamp"`
	Transactions  []Transaction `json:"transactions"`
	PrevHash      string        `json:"prev_hash"`
	Hash          string        `json:"hash"`
	Nonce         int64         `json:"nonce"`
	AI_Difficulty int           `json:"ai_difficulty"`
	Message       string        `json:"message,omitempty"`
}

var (
	Blockchain  []Block
	chainMutex  sync.Mutex
	historyFile = "neurocore_history.json"
)

func calculateHash(b Block) string {
	h := sha256.New()
	// simple deterministic serialization
	data, _ := json.Marshal(struct {
		Index      int
		Timestamp  int64
		TxCount    int
		PrevHash   string
		Nonce      int64
		Difficulty int
	}{b.Index, b.Timestamp, len(b.Transactions), b.PrevHash, b.Nonce, b.AI_Difficulty})
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func createGenesisBlock() Block {
	genesis := Block{
		Index:         0,
		Timestamp:     time.Now().Unix(),
		Transactions:  []Transaction{},
		PrevHash:      "0",
		Nonce:         0,
		AI_Difficulty: 2,
		Message:       "Neurocore: The First AI-Powered Ledger",
	}
	genesis.Hash = calculateHash(genesis)
	chainMutex.Lock()
	Blockchain = append(Blockchain, genesis)
	chainMutex.Unlock()
	persistChain()
	return genesis
}

func persistChain() error {
	chainMutex.Lock()
	defer chainMutex.Unlock()
	data, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(historyFile, data, 0644)
}

// GetAdaptiveDifficulty computes AI_Difficulty using the last blocks' timing
func GetAdaptiveDifficulty() int {
	chainMutex.Lock()
	defer chainMutex.Unlock()
	base := 2
	n := len(Blockchain)
	if n == 0 {
		return base
	}
	// every 5 blocks, adjust difficulty
	if n < 5 {
		return base
	}
	// compute average time of last 5 blocks
	end := n
	start := n - 5
	var sum float64
	for i := start; i < end; i++ {
		if i == 0 {
			continue
		}
		sum += float64(Blockchain[i].Timestamp - Blockchain[i-1].Timestamp)
	}
	avg := sum / 5.0
	// target seconds per block (simulated)
	target := 10.0
	diff := Blockchain[end-1].AI_Difficulty
	if diff == 0 {
		diff = base
	}
	if avg < target {
		diff++
	} else if avg > target*2 && diff > 1 {
		diff--
	}
	if diff < 1 {
		diff = 1
	}
	return diff
}

func AddBlock(b Block) error {
	chainMutex.Lock()
	defer chainMutex.Unlock()
	Blockchain = append(Blockchain, b)
	return persistChain()
}
