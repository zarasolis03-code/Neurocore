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
    PublicKey  []byte
    Address    string
}
func PubKeyToAddress(pub []byte) string {
    hash := sha256.Sum256(pub)
    return hex.EncodeToString(hash[:20])
}
func NewWallet() (*Wallet, error) {
    if _, err := os.Stat("wallet.key"); err == nil { return LoadWallet() }
    priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
    w := &Wallet{PrivateKey: priv, PublicKey: pub, Address: PubKeyToAddress(pub)}
    data, _ := x509.MarshalECPrivateKey(priv)
    os.WriteFile("wallet.key", data, 0600)
    return w, nil
}
func LoadWallet() (*Wallet, error) {
    data, _ := os.ReadFile("wallet.key")
    priv, _ := x509.ParseECPrivateKey(data)
    pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
    return &Wallet{PrivateKey: priv, PublicKey: pub, Address: PubKeyToAddress(pub)}, nil
}
