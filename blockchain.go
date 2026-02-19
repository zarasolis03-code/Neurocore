package main
import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
    "time"
)
type Block struct {
    Index        int
    Timestamp    string
    Transactions []string
    PrevHash     string
    Hash         string
    Nonce        int
}
func CalculateHash(block Block) string {
    record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, strings.Join(block.Transactions, ""), block.PrevHash, block.Nonce)
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}
func MineBlock(lastBlock Block, txs []string, diff int) Block {
    b := Block{lastBlock.Index + 1, time.Now().String(), txs, lastBlock.Hash, "", 0}
    prefix := strings.Repeat("0", diff)
    for {
        b.Hash = CalculateHash(b)
        if strings.HasPrefix(b.Hash, prefix) { break }
        b.Nonce++
    }
    return b
}
