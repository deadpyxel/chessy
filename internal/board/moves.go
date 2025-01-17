package board

import (
	"fmt"

	"github.com/deadpyxel/cheesy/internal/utils"
)

type MoveType uint8

// Types of chess movements
const (
	Normal    MoveType = iota // move a piece
	Capture                   // move to a square and capture opponent piece
	EnPassant                 // special pawn capture
	Castle                    // special movement between rook and king
	Promotion                 // special move for pawn, changing into another piece
)

// Move represents a chess move
type Move struct {
	From      Square   // starting position
	To        Square   // ending position
	Type      MoveType // type of chess move
	Promotion Piece    // Used for pawn promotion
}

// Container for possible moves
type MoveList struct {
	Moves [256]Move // provisional limit of number of moves generated
	Count int       // current amount of moves present on the MoveList
}

func (ml *MoveList) addMove(m Move) {
	if ml.Count < len(ml.Moves) {
		ml.Moves[ml.Count] = m
		ml.Count++
	}
}

func (ml *MoveList) String() string {
	var str string
	squares := make([]Square, ml.Count)
	for i := 0; i < ml.Count; i++ {
		str += ml.Moves[i].String() + "\n"
		squares[i] = ml.Moves[i].To
	}
	bb := Bitboard(0)
	for _, sq := range squares {
		bb = bb.Set(sq)
	}
	return fmt.Sprintf("\n%s\n%s", bb, str)
}

func (m Move) String() string {
	return m.From.String() + " -> " + m.To.String()
}

// Lookup table for movements
var (
	// Move positions for Knight and King cases
	KnightMoves = [8]int{-17, -15, -10, -6, 6, 10, 15, 17} // L shape movement (8 * rank + file), jumps over pieces
	KingMoves   = [8]int{-9, -8, -7, -1, 1, 7, 8, 9}       // only one square in any direction.

	// Directions for sliding pieces
	BishopDirections = [4]int{-9, -7, 7, 9}               // any amount of square in diagonals
	RookDirections   = [4]int{-8, -1, 1, 8}               // any amount of squares in file and ranks
	QueenDirections  = [8]int{-9, -8, -7, -1, 1, 7, 8, 9} // combined Bishop and Rook movements
)

func isOutOfBoard(sq int) bool {
	return sq < 0 || sq > 63
}

func (b *Board) generatePieceMoves(sq Square, piece Piece, color Color, ml *MoveList) {
	switch piece {
	case Pawn:
		b.generatePawnMoves(sq, color, ml)
	case Knight:
		b.generateKnightMoves(sq, color, ml)
	case Bishop:
		b.generateSlidingPieceMoves(sq, color, BishopDirections[:], ml)
	case Rook:
		b.generateSlidingPieceMoves(sq, color, RookDirections[:], ml)
	case Queen:
		b.generateSlidingPieceMoves(sq, color, QueenDirections[:], ml)
	case King:
		b.generateKingMoves(sq, color, ml)
	}
}

func (b *Board) generatePawnMoves(sq Square, color Color, ml *MoveList) {
	occOppColor := b.OccupiedByColor[color^1]
	occupied := b.OccupiedSquares

	fromFile := sq.FileOf()
	fromRank := sq.RankOf()

	var forward, startRank, prePromRank int // Determine direction and starting rank based on color
	if color == White {
		forward = 8     // direction of "forward" for pawns
		startRank = 1   // starting rank for pawns
		prePromRank = 6 // rank before promotion
	} else {
		forward = -8
		startRank = 6
		prePromRank = 1
	}

	// Normal movements
	singlePush := int(sq) + forward
	tgtSq := Square(singlePush)
	if !isOutOfBoard(singlePush) && !occupied.IsSet(tgtSq) {
		// we are promoting
		if fromRank == prePromRank {
			for piece := Knight; piece <= Queen; piece++ {
				ml.addMove(Move{
					From:      sq,
					To:        tgtSq,
					Type:      Promotion,
					Promotion: piece,
				})
			}
		} else {
			ml.addMove(Move{
				From: sq,
				To:   tgtSq,
				Type: Normal,
			})
		}

		// Double push (if on starting rank and single push is possible)
		if fromRank == startRank {
			doublePush := singlePush + forward
			tgtSq = Square(doublePush)
			if !occupied.IsSet(tgtSq) {
				ml.addMove(Move{
					From: sq,
					To:   tgtSq,
					Type: Normal,
				})
			}
		}
	}

	// Captures
	captureDir := []int{forward - 1, forward + 1} // front diagonals
	for _, dir := range captureDir {
		capSq := int(sq) + dir
		if isOutOfBoard(capSq) {
			continue
		}

		tgtSq = Square(capSq)
		tgtFile := tgtSq.FileOf()
		// check for edge wrapping
		if utils.Abs(fromFile-tgtFile) > 1 {
			continue
		}
		// Normal Captures
		if occOppColor.IsSet(tgtSq) {
			if fromRank == prePromRank {
				for piece := Knight; piece <= Queen; piece++ {
					ml.addMove(Move{
						From:      sq,
						To:        tgtSq,
						Type:      Capture | Promotion,
						Promotion: piece,
					})
				}
			} else {
				ml.addMove(Move{
					From: sq,
					To:   tgtSq,
					Type: Capture,
				})
			}
		}
	}
	// TODO: Add en passant captures
	// This will require:
	// 1. Tracking the last moved pawn
	// 2. Checking if it moved two squares
	// 3. Checking if our pawn is on the correct rank (5th for white, 4th for black)
}

func (b *Board) generateKingMoves(sq Square, color Color, ml *MoveList) {
	occSameColor := b.OccupiedByColor[color]  // squares occupied by same color
	occOppColor := b.OccupiedByColor[color^1] // squares ocuppied by opposite colors

	fromFile := sq.FileOf()

	for _, offset := range KingMoves {
		toSq := int(sq) + offset
		// Skip positions out of the board
		if isOutOfBoard(toSq) {
			continue
		}

		tgtSq := Square(toSq)
		toFile := tgtSq.FileOf()

		// Check if this move would wrap around the edges
		if utils.Abs(fromFile-toFile) > 1 {
			continue
		}

		if !occSameColor.IsSet(tgtSq) {
			mvType := Normal
			// Check if this is a capture
			if occOppColor.IsSet(tgtSq) {
				mvType = Capture
			}
			ml.addMove(Move{
				From: sq,
				To:   tgtSq,
				Type: mvType,
			})
		}
	}

	// TODO: Handle castling
	// it will need to track board state, if the king/rook has moved, if the path is clear
	// if king is not in check and if the squares related are not under attack
}

func (b *Board) generateKnightMoves(sq Square, color Color, ml *MoveList) {
	occSameColor := b.OccupiedByColor[color]  // squares occupied by same color
	occOppColor := b.OccupiedByColor[color^1] // squares ocuppied by opposite colors

	fromFile := sq.FileOf()

	for _, offset := range KnightMoves {
		toSq := int(sq) + offset
		// Skip out of bounds positions
		if isOutOfBoard(toSq) {
			continue
		}

		tgtSq := Square(toSq)
		toFile := tgtSq.FileOf()
		// Check if this move would wrap around board edges
		if utils.Abs(fromFile-toFile) > 2 {
			continue
		}
		// Skip cases where the square is occupied by the same color
		if !occSameColor.IsSet(tgtSq) {
			mvType := Normal
			// In case it is occupied by another color, mark it as a capture
			if occOppColor.IsSet(tgtSq) {
				mvType = Capture
			}

			ml.addMove(Move{
				From: sq,
				To:   tgtSq,
				Type: mvType,
			})
		}
	}
}

func (b *Board) generateSlidingPieceMoves(sq Square, color Color, directions []int, ml *MoveList) {
	occSameColor := b.OccupiedByColor[color]  // squares occupied by same color
	occOppColor := b.OccupiedByColor[color^1] // squares ocuppied by opposite colors
	fromFile := sq.FileOf()
	for _, dir := range directions {
		step := 1
		lastFile := fromFile
		for {
			toSq := int(sq) + (dir * step)
			// Out of bounds check
			if isOutOfBoard(toSq) {
				break
			}

			tgtSq := Square(toSq)
			toFile := tgtSq.FileOf()
			// Check if this move would wrap around board edges
			if utils.Abs(toFile-lastFile) > 1 {
				break
			}
			// Stop if we reach our own piece
			if occSameColor.IsSet(tgtSq) {
				break
			}

			mvType := Normal
			if occOppColor.IsSet(tgtSq) {
				mvType = Capture
				ml.addMove(Move{From: sq, To: tgtSq, Type: mvType})
				break
			}
			ml.addMove(Move{From: sq, To: tgtSq, Type: mvType})
			// Update current position
			lastFile = toFile
			step++
		}
	}
}
