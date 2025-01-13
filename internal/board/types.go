package board

// Custom Type for Piece
type Piece uint8

// Custom type for piece Color
type Color uint8

const (
	White Color = iota
	Black
	None
)

// Custom type for a Square in the board
type Square uint8

// Constants for Piece Enum
const (
	Empty Piece = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Bitboard uint64

// Set sets a bit in the Bitboard at the specified Square position.
func (bb Bitboard) Set(sq Square) Bitboard {
	return bb | (1 << sq)
}

// Clear sets the bit at the given square to 0 in the Bitboard and returns the updated Bitboard.
func (bb Bitboard) Clear(sq Square) Bitboard {
	return bb &^ (1 << sq)
}

// IsSet checks if a specific bit is set in the given Bitboard.
func (bb Bitboard) IsSet(sq Square) bool {
	return (bb & (1 << sq)) != 0
}

func (bb Bitboard) String() string {
	var result string
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sq := Square(rank*8 + file)
			if bb.IsSet(sq) {
				result += "1 "
			} else {
				result += "0 "
			}
		}
		result += "\n"
	}
	return result
}

const (
	FileA Bitboard = 0x0101010101010101
	FileB Bitboard = FileA << 1
	FileC Bitboard = FileA << 2
	FileD Bitboard = FileA << 3
	FileE Bitboard = FileA << 4
	FileF Bitboard = FileA << 5
	FileG Bitboard = FileA << 6
	FileH Bitboard = FileA << 7

	Rank1 Bitboard = 0x00000000000000FF
	Rank2 Bitboard = Rank1 << (8 * 1)
	Rank3 Bitboard = Rank1 << (8 * 2)
	Rank4 Bitboard = Rank1 << (8 * 3)
	Rank5 Bitboard = Rank1 << (8 * 4)
	Rank6 Bitboard = Rank1 << (8 * 5)
	Rank7 Bitboard = Rank1 << (8 * 6)
	Rank8 Bitboard = Rank1 << (8 * 7)
)

type Board struct {
	Pieces [2][7]Bitboard // Bitboards for each piecetype [Color][Type]

	// Combined Bitboards for faster lookup
	OccupiedSquares Bitboard    // All pieces
	OccupiedByColor [2]Bitboard // All pieces of the same color

	// Positional information
	SideToMove Color
}

// SetInitialBoard initializes the chess board with the starting positions of all pieces.
func (b *Board) SetInitialBoard() {
	// White pieces
	b.Pieces[White][Pawn] = Rank2
	b.Pieces[White][Rook] = (1 << 0) | (1 << 7)   // A1 and H1
	b.Pieces[White][Knight] = (1 << 1) | (1 << 6) // B1 and G1
	b.Pieces[White][Bishop] = (1 << 2) | (1 << 5) // C1 and F1
	b.Pieces[White][Queen] = (1 << 3)             // D1
	b.Pieces[White][King] = (1 << 4)              // E1
	// Black pieces (White pieces shifted by 56 squares)
	b.Pieces[Black][Pawn] = Rank7
	b.Pieces[Black][Rook] = b.Pieces[White][Rook] << 56
	b.Pieces[Black][Knight] = b.Pieces[White][Knight] << 56
	b.Pieces[Black][Bishop] = b.Pieces[White][Bishop] << 56
	b.Pieces[Black][Queen] = b.Pieces[White][Queen] << 56
	b.Pieces[Black][King] = b.Pieces[White][King] << 56

	// Update all occupied squares on the boards
	b.UpdateOccupiedSquares()
}

// UpdateOccupiedSquares updates the occupied squares on the board for both white and black pieces.
func (b *Board) UpdateOccupiedSquares() {
	b.OccupiedByColor[White] = 0
	b.OccupiedByColor[Black] = 0

	for piece := Pawn; piece <= King; piece++ {
		b.OccupiedByColor[White] |= b.Pieces[White][piece]
		b.OccupiedByColor[Black] |= b.Pieces[Black][piece]
	}

	b.OccupiedSquares = b.OccupiedByColor[White] | b.OccupiedByColor[Black]
}

// GetPieceAt returns the piece and its color at the given square
func (b *Board) GetPieceAt(sq Square) (Color, Piece) {
	for color := White; color <= Black; color++ {
		for piece := Pawn; piece <= King; piece++ {
			if b.Pieces[color][piece].IsSet(sq) {
				return color, piece
			}
		}
	}
	return None, Empty
}
