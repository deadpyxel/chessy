package board

// Custom Type for Piece
type Piece uint8

// Custom type for piece Color
type Color uint8

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
