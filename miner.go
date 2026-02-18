package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

var globalHashCounter uint64

// Helper to count leading zero bits in a hex-encoded hash
func leadingZeroBits(hexHash string) int {
	b, _ := hex.DecodeString(hexHash)
	bits := 0
	for _, by := range b {
		for i := 7; i >= 0; i-- {
			if (by>>i)&1 == 0 {
				bits++
			} else {
				return bits
			}
		}
	}
	return bits
}

// MatrixHash simulates an AI workload; returns hex string of a sha256 of small matrix ops
func MatrixHash(seed string, nonce int64) string {
	// create a small pseudo-matrix from seed and nonce
	h := sha256.New()
	h.Write([]byte(seed))
	nBytes := []byte(fmt.Sprintf("%d", nonce))
	h.Write(nBytes)
	base := h.Sum(nil)

	// small simulated matrix multiply loop
	s := int(math.Min(16, float64(len(base))))
	mat := make([]int, s)
	for i := 0; i < s; i++ {
		mat[i] = int(base[i])
	}
	acc := 0
	// pseudo-compute to simulate CPU work
	for i := 0; i < 128; i++ {
		for j := 0; j < s; j++ {
			acc = (acc*mat[j] + j + int(nonce) + i) % 0xFFFF
		}
	}
	final := sha256.Sum256(append(base, byte(acc&0xFF)))
	return hex.EncodeToString(final[:])
}

// Miner consumes a mempool snapshot and tries nonces until difficulty met.
func Miner(mempoolCh <-chan Transaction, newBlockCh chan<- Block, stopCh <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-stopCh:
			return
		default:
			// assemble transactions: drain mempool channel non-blocking
			txs := []Transaction{}
		gather:
			for {
				select {
				case tx := <-mempoolCh:
					txs = append(txs, tx)
				default:
					break gather
				}
			}

			// build candidate block
			chainLen := 0
			var prevHash string
			var index int
			chainLen = len(Blockchain)
			if chainLen > 0 {
				prev := Blockchain[chainLen-1]
				prevHash = prev.Hash
				index = prev.Index + 1
			} else {
				prevHash = "0"
				index = 0
			}
			difficulty := GetAdaptiveDifficulty()

			block := Block{
				Index:         index,
				Timestamp:     time.Now().Unix(),
				Transactions:  txs,
				PrevHash:      prevHash,
				Nonce:         0,
				AI_Difficulty: difficulty,
			}

			// mining loop
			start := time.Now()
			var nonce int64 = 0
			for {
				// simulate some non-deterministic work per try
				nonce = int64(rand.Intn(1 << 30))
				seed := fmt.Sprintf("%d:%s:%d", block.Index, block.PrevHash, time.Now().UnixNano())
				h := MatrixHash(seed, nonce)
				atomic.AddUint64(&globalHashCounter, 1)
				if leadingZeroBits(h) >= block.AI_Difficulty {
					block.Nonce = nonce
					block.Timestamp = time.Now().Unix()
					block.Hash = h
					// attach message for genesis-like transparency
					if block.Index == 0 {
						block.Message = "Neurocore: The First AI-Powered Ledger"
					}
					// broadcast (send back)
					newBlockCh <- block
					break
				}
				// throttle to yield occasionally
				if time.Since(start) > 5*time.Second {
					select {
					case <-ticker.C:
					default:
					}
				}
				// continue trying
			}
		}
	}
}
