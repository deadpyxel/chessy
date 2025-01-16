package board

import (
	"fmt"
	"testing"
)

func extractMoves(ml *MoveList) map[Square]Move {
	genMoves := make(map[Square]Move)
	for i := 0; i < ml.Count; i++ {
		move := ml.Moves[i]
		genMoves[move.To] = move
	}
	return genMoves
}

func TestIsOutOfBounds(t *testing.T) {
	tests := []struct {
		name     string
		sq       int
		expected bool
	}{
		{name: "when under 0 returns true", sq: -1, expected: true},
		{name: "when above 63 returns true", sq: 64, expected: true},
		{name: "when between 0 and 63 returns false", sq: 1, expected: false},
		{name: "when is 0 returns false", sq: 0, expected: false},
		{name: "when is 63 returns false", sq: 63, expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isOutOfBounds(test.sq)
			if actual != test.expected {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*Board)
		startSq       Square
		expectedMoves []Move
	}{
		{
			name:    "white knight in center of empty board",
			startSq: Square(28), // e4
			setup: func(b *Board) {
				b.Pieces[White][Knight] = Bitboard(1 << 28)
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				{From: 28, To: 11, Type: Normal}, // e4 to c3
				{From: 28, To: 13, Type: Normal}, // e4 to f3
				{From: 28, To: 18, Type: Normal}, // e4 to d2
				{From: 28, To: 22, Type: Normal}, // e4 to f2
				{From: 28, To: 34, Type: Normal}, // e4 to d6
				{From: 28, To: 38, Type: Normal}, // e4 to f6
				{From: 28, To: 43, Type: Normal}, // e4 to c5
				{From: 28, To: 45, Type: Normal}, // e4 to f5
			},
		},
		{
			name:    "white knight in corner (a1)",
			startSq: Square(0), // a1
			setup: func(b *Board) {
				b.Pieces[White][Knight] = Bitboard(1 << 0)
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				{From: 0, To: 10, Type: Normal}, // a1 to c2
				{From: 0, To: 17, Type: Normal}, // a1 to b3
			},
		},
		{
			name:    "white knight with capture opportunities",
			startSq: Square(28), // e4
			setup: func(b *Board) {
				b.Pieces[White][Knight] = Bitboard(1 << 28)
				b.Pieces[Black][Pawn] = Bitboard(1<<11 | 1<<38) // Place black pawns on c3 and f6
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				{From: 28, To: 11, Type: Capture}, // e4 to c3 (capture)
				{From: 28, To: 13, Type: Normal},  // e4 to f3
				{From: 28, To: 18, Type: Normal},  // e4 to d2
				{From: 28, To: 22, Type: Normal},  // e4 to f2
				{From: 28, To: 34, Type: Normal},  // e4 to d6
				{From: 28, To: 38, Type: Capture}, // e4 to f6 (capture)
				{From: 28, To: 43, Type: Normal},  // e4 to c5
				{From: 28, To: 45, Type: Normal},  // e4 to f5
			},
		},
		{
			name:    "white knight blocked by friendly pieces",
			startSq: Square(28), // e4
			setup: func(b *Board) {
				b.Pieces[White][Knight] = Bitboard(1 << 28)
				b.Pieces[White][Pawn] = Bitboard(1<<11 | 1<<38) // Place white pawns on c3 and f6
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				{From: 28, To: 13, Type: Normal}, // e4 to f3
				{From: 28, To: 18, Type: Normal}, // e4 to d2
				{From: 28, To: 22, Type: Normal}, // e4 to f2
				{From: 28, To: 34, Type: Normal}, // e4 to d6
				{From: 28, To: 43, Type: Normal}, // e4 to c5
				{From: 28, To: 45, Type: Normal}, // e4 to f5
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)

			var ml MoveList
			b.generateKnightMoves(tt.startSq, White, &ml)
			if ml.Count != len(tt.expectedMoves) {
				t.Errorf("Expected moveset to have %d entries, got %d instead", len(tt.expectedMoves), ml.Count)
				fmt.Printf("%v\n", ml.Moves[:ml.Count])
			}

			genMoves := extractMoves(&ml)
			for _, expMove := range tt.expectedMoves {
				move, exists := genMoves[expMove.To]
				if !exists {
					t.Errorf("Expected move to %s not found", expMove.To)
					continue
				}
				if move.Type != expMove.Type {
					t.Errorf("Expected move to %s type to be %d, got %d instead", expMove.To, expMove.Type, move.Type)
				}
			}
		})
	}
}

func TestGenerateSlidingMoves(t *testing.T) {
	tests := []struct {
		name          string
		piece         Piece
		setup         func(*Board)
		startSq       Square
		expectedMoves []Move
	}{
		{
			name:    "bishop in e4 center of empty board",
			piece:   Bishop,
			startSq: Square(28), // e4
			setup: func(b *Board) {
				b.Pieces[White][Bishop] = Bitboard(1 << 28)
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				// Northeast diagonal
				{From: 28, To: 37, Type: Normal},
				{From: 28, To: 46, Type: Normal},
				{From: 28, To: 55, Type: Normal},
				// Southeast diagonal
				{From: 28, To: 19, Type: Normal},
				{From: 28, To: 10, Type: Normal},
				{From: 28, To: 1, Type: Normal},
				// Southwest diagonal
				{From: 28, To: 21, Type: Normal},
				{From: 28, To: 14, Type: Normal},
				{From: 28, To: 7, Type: Normal},
				// Northwest diagonal
				{From: 28, To: 35, Type: Normal},
				{From: 28, To: 42, Type: Normal},
				{From: 28, To: 49, Type: Normal},
			},
		},
		{
			name:    "rook in e4 with capture opportunity and blocked by friendly piece",
			piece:   Rook,
			startSq: Square(28), // e4
			setup: func(b *Board) {
				b.Pieces[White][Rook] = Bitboard(1 << 28)
				b.Pieces[White][Pawn] = Bitboard(1 << 36) // Friendly pawn on e5
				b.Pieces[Black][Pawn] = Bitboard(1 << 20) // Enemy pawn on e3
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				// North (blocked by friendly pawn)
				// South (until enemy pawn)
				{From: 28, To: 20, Type: Capture},
				// East
				{From: 28, To: 29, Type: Normal},
				{From: 28, To: 30, Type: Normal},
				{From: 28, To: 31, Type: Normal},
				// West
				{From: 28, To: 27, Type: Normal},
				{From: 28, To: 26, Type: Normal},
				{From: 28, To: 25, Type: Normal},
				{From: 28, To: 24, Type: Normal},
			},
		},
		{
			name:    "queen in corner a1 with mixed blocking",
			piece:   Queen,
			startSq: Square(0), // a1
			setup: func(b *Board) {
				b.Pieces[White][Queen] = Bitboard(1 << 0)
				b.Pieces[Black][Pawn] = Bitboard(1 << 9) // Enemy pawn on b2
				b.UpdateOccupiedSquares()
			},
			expectedMoves: []Move{
				// Diagonal (northeast only from a1)
				{From: 0, To: 9, Type: Capture},
				// Vertical (north from a1)
				{From: 0, To: 8, Type: Normal},
				{From: 0, To: 16, Type: Normal},
				{From: 0, To: 24, Type: Normal},
				{From: 0, To: 32, Type: Normal},
				{From: 0, To: 40, Type: Normal},
				{From: 0, To: 48, Type: Normal},
				{From: 0, To: 56, Type: Normal},
				// Horizontal (east from a1)
				{From: 0, To: 1, Type: Normal},
				{From: 0, To: 2, Type: Normal},
				{From: 0, To: 3, Type: Normal},
				{From: 0, To: 4, Type: Normal},
				{From: 0, To: 5, Type: Normal},
				{From: 0, To: 6, Type: Normal},
				{From: 0, To: 7, Type: Normal},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)

			var ml MoveList
			var directions []int
			switch tt.piece {
			case Bishop:
				directions = BishopDirections[:]
			case Rook:
				directions = RookDirections[:]
			case Queen:
				directions = QueenDirections[:]
			}

			b.generateSlidingPieceMoves(tt.startSq, White, directions, &ml)

			if ml.Count != len(tt.expectedMoves) {
				t.Errorf("Expected moveset to have %d entries, got %d instead", len(tt.expectedMoves), ml.Count)
				fmt.Printf("%v\n", ml.Moves[:ml.Count])
			}

			genMoves := extractMoves(&ml)
			for _, expMove := range tt.expectedMoves {
				move, exists := genMoves[expMove.To]
				if !exists {
					t.Errorf("Expected move to %s not found", expMove.To)
					continue
				}
				if move.Type != expMove.Type {
					t.Errorf("Expected move to %s type to be %d, got %d instead", expMove.To, expMove.Type, move.Type)
				}
			}
		})
	}
}
