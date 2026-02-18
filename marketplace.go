package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Listing struct {
	ID        string  `json:"id"`
	Seller    string  `json:"seller"`
	Asset     string  `json:"asset"`
	Price     float64 `json:"price"`
	CreatedAt int64   `json:"created_at"`
	Sold      bool    `json:"sold"`
}

type Marketplace struct {
	mu       sync.Mutex
	listings map[string]*Listing
}

var Market = &Marketplace{listings: make(map[string]*Listing)}

func (m *Marketplace) CreateListing(seller, asset string, price float64) *Listing {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := generateID()
	l := &Listing{ID: id, Seller: seller, Asset: asset, Price: price, CreatedAt: time.Now().Unix(), Sold: false}
	m.listings[id] = l
	return l
}

func (m *Marketplace) ListListings() []*Listing {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Listing, 0, len(m.listings))
	for _, l := range m.listings {
		out = append(out, l)
	}
	return out
}

func (m *Marketplace) Buy(listingID, buyer string) (*Listing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.listings[listingID]
	if !ok {
		return nil, errors.New("listing not found")
	}
	if l.Sold {
		return nil, errors.New("already sold")
	}
	// mark sold and in a real chain create transfer tx
	l.Sold = true
	// optionally create transaction to transfer asset (simulated)
	tx := Transaction{From: l.Seller, To: buyer, Amount: l.Price}
	// push into mempool if available
	select {
	case mempoolCh <- tx:
	default:
	}
	return l, nil
}

func generateID() string {
	// simple unique id for listings
	return fmtNowUnixString()
}

func fmtNowUnixString() string {
	return fmt.Sprintf("L-%d", time.Now().UnixNano())
}
