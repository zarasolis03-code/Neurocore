package main

import (
	"testing"
)

func TestPubKeyToAddressLength(t *testing.T) {
	w, err := NewWallet()
	if err != nil {
		t.Fatal(err)
	}
	if len(w.Address) == 0 {
		t.Fatalf("empty address")
	}
}
