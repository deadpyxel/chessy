package board

import (
	"fmt"
	"testing"
)

func TestBitboardSet(t *testing.T) {
	var bb Bitboard = 0
	tests := []struct {
		name     string
		initial  Bitboard
		sq       Square
		expected Bitboard
	}{
		{
			name:     "set a1 on empty board",
			initial:  bb,
			sq:       Square(0),
			expected: Bitboard(0x1),
		},
		{
			name:     "set h8 on empty board",
			initial:  bb,
			sq:       Square(63),
			expected: Bitboard(0x8000000000000000),
		},
		{
			name:     "set square already set",
			initial:  Bitboard(0x1),
			sq:       Square(0),
			expected: Bitboard(0x1),
		},
		{
			name:     "set square on full board set",
			initial:  Bitboard(0xFFFFFFFFFFFFFFFF),
			sq:       Square(31),
			expected: Bitboard(0xFFFFFFFFFFFFFFFF),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.initial.Set(tt.sq)
			if res != tt.expected {
				fmt.Println(res)
				t.Errorf("Expected bitboard to have bit set at square %d, but it doesn't", tt.sq)
			}
		})
	}
}

func TestBitboardClear(t *testing.T) {
	var bb Bitboard = 0xFFFFFFFFFFFFFFFF
	tests := []struct {
		name     string
		initial  Bitboard
		sq       Square
		expected Bitboard
	}{
		{
			name:     "clear a1 on full board",
			initial:  bb,
			sq:       Square(0),
			expected: Bitboard(0xFFFFFFFFFFFFFFFE),
		},
		{
			name:     "clear h8 on full board",
			initial:  bb,
			sq:       Square(63),
			expected: Bitboard(0x7FFFFFFFFFFFFFFF),
		},
		{
			name:     "clear square on empty board",
			initial:  Bitboard(0),
			sq:       Square(63),
			expected: Bitboard(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.initial.Clear(tt.sq)
			if res != tt.expected {
				fmt.Println(res)
				t.Errorf("Expected bitboard to have bit cleared at square %d, but it is still set", tt.sq)
			}
		})
	}
}

func TestBitboardIsSet(t *testing.T) {
	tests := []struct {
		name     string
		initial  Bitboard
		sq       Square
		expected bool
	}{
		{
			name:     "checking if square is set on empty board returns false",
			initial:  Bitboard(0),
			sq:       Square(0),
			expected: false,
		},
		{
			name:     "checking if square is set on full board returns true",
			initial:  Bitboard(0xFFFFFFFFFFFFFFFF),
			sq:       Square(0),
			expected: true,
		},
		{
			name:     "checking if set square is set returns true",
			initial:  Bitboard(0x1),
			sq:       Square(0),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.initial.IsSet(tt.sq)
			if res != tt.expected {
				fmt.Println(tt.initial)
				t.Errorf("Expected bitboard square %d to be set to %v, got %v instead", tt.sq, tt.expected, res)
			}
		})
	}
}
