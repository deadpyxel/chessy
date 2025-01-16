package board

import "github.com/deadpyxel/cheesy/internal/utils"

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

func isOutOfBounds(sq int) bool {
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

}

func (b *Board) generateKingMoves(sq Square, color Color, ml *MoveList) {

}

func (b *Board) generateKnightMoves(sq Square, color Color, ml *MoveList) {
	occSameColor := b.OccupiedByColor[color]  // squares occupied by same color
	occOppColor := b.OccupiedByColor[color^1] // squares ocuppied by opposite colors

	fromFile := sq.FileOf()

	for _, offset := range KnightMoves {
		toSq := int(sq) + offset
		// Skip out of bounds positions
		if isOutOfBounds(toSq) {
			continue
		}

		tgtSq := Square(toSq)
		toFile := tgtSq.FileOf()
		// Check if this move would wrap around board edges
		if fileDiff := utils.Abs(fromFile - toFile); fileDiff > 2 {
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
	fromRank := sq.RankOf()
	for _, dir := range directions {
		toSq := int(sq)
		currFile := fromFile // Initialize currFile for this direction
		currRank := fromRank

		for {
			toSq += dir
			// Out of bounds check
			if isOutOfBounds(toSq) {
				break
			}

			tgtSq := Square(toSq)
			toFile := tgtSq.FileOf()
			toRank := tgtSq.RankOf()

			fileDiff := utils.Abs(toFile - currFile)
			rankDiff := utils.Abs(toRank - currRank)
			absDir := utils.Abs(dir)

			// For diagonal moves the difference should be only of 1
			isDiagonal := absDir == 7 || absDir == 9
			if isDiagonal && (fileDiff != 1 || rankDiff != 1) {
				break
			}
			// For horizontal moves only the file should change
			isHorizontal := absDir == 1
			if isHorizontal && rankDiff != 0 {
				break
			}
			// For vertical moves only the rank should change
			isVertical := absDir == 8
			if isVertical && fileDiff != 0 {
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
			currFile = toFile
			currRank = toRank
		}
	}
}
