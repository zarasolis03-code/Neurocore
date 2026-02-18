package main
import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "time"
)
type Transaction struct { From, To string; Amount float64; Timestamp int64 }
type Block struct {
    Index int; Timestamp int64; PrevHash, Hash string; Nonce int64
    AI_Difficulty int; Miner string
}
var (
    Blockchain []Block
    mu sync.Mutex
)
func calculateHash(b Block) string {
    data := fmt.Sprintf("%d%d%s%d%d", b.Index, b.Timestamp, b.PrevHash, b.Nonce, b.AI_Difficulty)
    h := sha256.Sum256([]byte(data))
    return hex.EncodeToString(h[:])
}
func LoadChain() {
    mu.Lock(); defer mu.Unlock()
    if _, err := os.Stat("blockchain.json"); os.IsNotExist(err) {
        gen := Block{Index: 0, Timestamp: time.Now().Unix(), PrevHash: "0", AI_Difficulty: 2, Miner: "Genesis"}
        gen.Hash = calculateHash(gen)
        Blockchain = append(Blockchain, gen)
    } else {
        data, _ := os.ReadFile("blockchain.json")
        json.Unmarshal(data, &Blockchain)
    }
}
func SaveChain() {
    data, _ := json.MarshalIndent(Blockchain, "", "  ")
    os.WriteFile("blockchain.json", data, 0644)
}
func AddBlock(b Block) { mu.Lock(); defer mu.Unlock(); Blockchain = append(Blockchain, b); SaveChain() }
func GetLatestBlock() Block { mu.Lock(); defer mu.Unlock(); return Blockchain[len(Blockchain)-1] }
func GetBalance(addr string) float64 {
    mu.Lock(); defer mu.Unlock(); bal := 0.0
    for _, b := range Blockchain { if b.Miner == addr { bal += 50 } }
    return bal
}
