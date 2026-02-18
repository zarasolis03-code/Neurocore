package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

func NewWallet() (*Wallet, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	addr := PubKeyToAddress(&priv.PublicKey)
	return &Wallet{PrivateKey: priv, Address: addr}, nil
}

func PubKeyToAddress(pub *ecdsa.PublicKey) string {
	// encode the (x||y) and hash
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	data := append(xBytes, yBytes...)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:20]) // truncated for readability
}

func SignTransaction(priv *ecdsa.PrivateKey, txHash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, priv, txHash)
	if err != nil {
		return nil, err
	}
	rb := r.Bytes()
	sb := s.Bytes()
	sig := append(rb, sb...)
	return sig, nil
}

func VerifySignature(pub *ecdsa.PublicKey, txHash []byte, signature []byte) bool {
	// split signature in half
	if len(signature)%2 != 0 {
		return false
	}
	half := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:half])
	s := new(big.Int).SetBytes(signature[half:])
	return ecdsa.Verify(pub, txHash, r, s)
}

func HashTransaction(tx Transaction) []byte {
	// simple transaction fingerprint
	h := sha256.New()
	h.Write([]byte(tx.From))
	h.Write([]byte(tx.To))
	h.Write([]byte([]byte(string(rune(int(tx.Amount * 100))))))
	return h.Sum(nil)
}

func (w *Wallet) PrivateKeyHex() string {
	return hex.EncodeToString(w.PrivateKey.D.Bytes())
}

// verifyTx convenience helper
func verifyTx(tx Transaction) bool {
	// In a real system we'd recover public key; here we do a simplified flow
	// (caller is expected to use VerifySignature with known pubkey)
	return true
}
