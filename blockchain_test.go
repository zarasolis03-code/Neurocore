package main

import "testing"

func TestAdaptiveDifficultyBounds(t *testing.T) {
	// reset chain
	Blockchain = []Block{}
	// create 6 blocks with timestamps spaced to force difficulty change
	for i := 0; i < 6; i++ {
		b := Block{Index: i, Timestamp: int64(100 + i*5), AI_Difficulty: 2}
		Blockchain = append(Blockchain, b)
	}
	d := GetAdaptiveDifficulty()
	if d < 1 {
		t.Fatalf("difficulty too low: %d", d)
	}
}
