package main
import (
    "time"
)
func Miner(newBlockCh chan Block, stopCh chan struct{}) {
    w, _ := NewWallet()
    for {
        select {
        case <-stopCh: return
        default:
            last := GetLatestBlock()
            diff := 3
            nb := Block{Index: last.Index+1, Timestamp: time.Now().Unix(), PrevHash: last.Hash, AI_Difficulty: diff, Miner: w.Address}
            for i := int64(0); ; i++ {
                nb.Nonce = i; nb.Hash = calculateHash(nb)
                if nb.Hash[:diff] == "000" {
                    newBlockCh <- nb
                    break
                }
            }
        }
    }
}
