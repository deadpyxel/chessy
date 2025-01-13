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
