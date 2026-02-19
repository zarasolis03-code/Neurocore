package main
import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "crypto/x509"
    "encoding/hex"
    "os"
)
type Wallet struct {
    PrivateKey *ecdsa.PrivateKey
    Address    string
}
func NewWallet() (*Wallet, error) {
    if _, err := os.Stat("wallet.key"); err == nil {
        data, _ := os.ReadFile("wallet.key")
        priv, _ := x509.ParseECPrivateKey(data)
        pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
        hash := sha256.Sum256(pub)
        return &Wallet{priv, hex.EncodeToString(hash[:20])}, nil
    }
    priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
    hash := sha256.Sum256(pub)
    data, _ := x509.MarshalECPrivateKey(priv)
    os.WriteFile("wallet.key", data, 0600)
    return &Wallet{priv, hex.EncodeToString(hash[:20])}, nil
}
