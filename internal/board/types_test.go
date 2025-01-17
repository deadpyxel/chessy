package board

import (
	"fmt"
	"testing"
)

func TestSquareFileOf(t *testing.T) {
	tests := []struct {
		sq       Square
		wantFile int
	}{
		{sq: Square(0), wantFile: 0},  // file a
		{sq: Square(63), wantFile: 7}, // file h
		{sq: Square(28), wantFile: 4}, // file e
		{sq: Square(10), wantFile: 2}, // file c
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Square %s FileOf() returns %d", tt.sq, tt.wantFile), func(t *testing.T) {
			result := tt.sq.FileOf()
			if result != tt.wantFile {
				t.Errorf("Expected %d, but got %d", tt.wantFile, result)
			}
		})
	}
}

func TestSquareRankOf(t *testing.T) {
	tests := []struct {
		sq       Square
		wantRank int
	}{
		{sq: Square(0), wantRank: 0},  // rank 1
		{sq: Square(63), wantRank: 7}, // rank 8
		{sq: Square(24), wantRank: 3}, // rank 4
		{sq: Square(10), wantRank: 1}, // rank 2
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Square %s RankOf() returns %d", tt.sq, tt.wantRank), func(t *testing.T) {
			result := tt.sq.RankOf()
			if result != tt.wantRank {
				t.Errorf("Expected %d, but got %d", tt.wantRank, result)
			}
		})
	}
}

func TestSquaresSweepingNotation(t *testing.T) {
	for sq := Square(0); sq < 64; sq++ {
		notation := sq.String()

		// Check length
		if len(notation) != 2 {
			t.Errorf("Square(%d).String() produced invalid length notation: %s", sq, notation)
			continue
		}

		// Check file character
		if notation[0] < 'a' || notation[0] > 'h' {
			t.Errorf("Square(%d).String() produced invalid file character: %c", sq, notation[0])
		}

		// Check rank character
		if notation[1] < '1' || notation[1] > '8' {
			t.Errorf("Square(%d).String() produced invalid rank character: %c", sq, notation[1])
		}
	}
}
