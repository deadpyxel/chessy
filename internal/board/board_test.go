package board

import (
	"testing"
)

func TestBoardGetPieceAt(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Board)
		sq        Square
		wantColor Color
		wantPiece Piece
	}{
		{
			name:      "lookup on empty board returns no color and empty space",
			setup:     func(b *Board) {},
			sq:        Square(0),
			wantColor: None,
			wantPiece: Empty,
		},
		{
			name: "lookup on specific square with piece board returns matching color and type",
			setup: func(b *Board) {
				b.Pieces[White][Rook] = Rank1 & FileA
			},
			sq:        Square(0),
			wantColor: White,
			wantPiece: Rook,
		},
		{
			name: "lookup on specific square with 2 piece board returns matching color and type",
			setup: func(b *Board) {
				b.Pieces[White][Rook] = (Rank1 & FileA)
				b.Pieces[Black][Rook] = (Rank1 & FileH)
			},
			sq:        Square(7),
			wantColor: Black,
			wantPiece: Rook,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)
			gotColor, gotPiece := b.GetPieceAt(tt.sq)
			if gotColor != tt.wantColor {
				t.Errorf("Board.GetPieceAt() gotColor = %v, want %v", gotColor, tt.wantColor)
			}
			if gotPiece != tt.wantPiece {
				t.Errorf("Board.GetPieceAt() gotPiece = %v, want %v", gotPiece, tt.wantPiece)
			}
		})
	}
}

func TestBoardUpdateOccupiedSquares(t *testing.T) {
	b := &Board{}
	b.SetInitialBoard()

	// Test initial occupied squares
	expectedOccupiedByColor := map[Color]Bitboard{
		White: b.Pieces[White][Pawn] | b.Pieces[White][Knight] | b.Pieces[White][Bishop] | b.Pieces[White][Rook] | b.Pieces[White][Queen] | b.Pieces[White][King],
		Black: b.Pieces[Black][Pawn] | b.Pieces[Black][Knight] | b.Pieces[Black][Bishop] | b.Pieces[Black][Rook] | b.Pieces[Black][Queen] | b.Pieces[Black][King],
	}
	expectedOccupiedSquares := expectedOccupiedByColor[White] | expectedOccupiedByColor[Black]

	b.UpdateOccupiedSquares()

	if expectedOccupiedByColor[White] != b.OccupiedByColor[White] {
		t.Errorf("Expected all White pieces to occupy the following squares:\n%s\n got \n%s\n", b.OccupiedByColor[White], expectedOccupiedByColor[White])
	}

	if expectedOccupiedByColor[Black] != b.OccupiedByColor[Black] {
		t.Errorf("Expected all Black pieces to occupy the following squares:\n%s\n got \n%s\n", b.OccupiedByColor[Black], expectedOccupiedByColor[Black])
	}
	if expectedOccupiedSquares != b.OccupiedSquares {
		t.Errorf("Expected all pieces to occupy the following squares:\n%s\n got \n%s\n", b.OccupiedSquares, expectedOccupiedSquares)
	}
}

func TestPlayMoveNormal(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*Board)
		mv            Move
		expectedBoard func(*Board)
	}{
		{
			name: "white pawn on e4 can move to e5",
			setup: func(b *Board) {
				b.Pieces[White][Pawn] = Bitboard(1 << 28)
				b.SideToMove = White
				b.UpdateOccupiedSquares()
			},
			mv: Move{From: Square(28), To: Square(36), Type: Normal},
			expectedBoard: func(b *Board) {
				b.Pieces[White][Pawn] = Bitboard(1 << 36)
				b.SideToMove = Black
				b.UpdateOccupiedSquares()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)

			err := b.PlayMove(tt.mv)
			if err != nil {
				t.Errorf("Expected no error, got %v instead", err)
			}

			expectedBoard := &Board{}
			tt.expectedBoard(expectedBoard)

			// After the move is played, the active playing side should change
			if b.SideToMove != expectedBoard.SideToMove {
				t.Errorf("side to move should be to %v after playing the move, got %v instead", expectedBoard.SideToMove, b.SideToMove)
			}

			if !b.isEqualBoard(*expectedBoard) {
				t.Errorf("resulting board does not match desired output")
			}
		})
	}
}

func TestBoardToFEN(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Board)
		expected string
	}{
		{
			name: "empty board",
			setup: func(b *Board) {
				// Empty
				b.FullMoveCount = 1 // particular case where we have the move counter as zero

			},
			expected: "8/8/8/8/8/8/8/8 w KQkq - 0 1",
		},
		{
			name: "initial position",
			setup: func(b *Board) {
				b.SetInitialBoard()
			},
			expected: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "after 1. e4",
			setup: func(b *Board) {
				b.SetInitialBoard()
				b.PlayMove(Move{
					From: Square(12), // e2
					To:   Square(28), // e4
					Type: Normal,
				})
			},
			expected: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name: "complex middle game position",
			setup: func(b *Board) {
				b.Pieces[White][Pawn] = Bitboard(1<<18 | 1<<20 | 1<<21 | 1<<28) // c3, e3, f3, e4
				b.Pieces[White][Knight] = Bitboard(1 << 42)                     // Nc6
				b.Pieces[White][King] = Bitboard(1 << 6)                        // Kg1
				b.Pieces[Black][Pawn] = Bitboard(1<<49 | 1<<50 | 1<<51 | 1<<52) // b7, c7, d7, e7
				b.Pieces[Black][Queen] = Bitboard(1 << 40)                      // Qa6
				b.FullMoveCount = 1
				b.UpdateOccupiedSquares()
			},
			expected: "8/1pppp3/q1N5/8/4P3/2P1PP2/8/6K1 w KQkq - 0 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Board{}
			tt.setup(b)

			result := b.ToFEN()
			if result != tt.expected {
				t.Errorf("expected FEN to be %s, got %s instead", tt.expected, result)
			}
		})
	}
}
