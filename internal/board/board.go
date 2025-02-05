package board

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Board) PlayMove(m Move) error {
	// get moving piece from the board
	pCol, piece := b.GetPieceAt(m.From)
	if pCol == None || piece == Empty {
		return fmt.Errorf("no piece at source square %v", m.From)
	}
	// TODO: add checks for trying to move a piece not owned
	if pCol != b.SideToMove {
		return fmt.Errorf("cannot move opponent piece")
	}

	// Handle different move types
	switch m.Type {
	case Normal:
		b.movePiece(pCol, piece, m.From, m.To)
	case Capture:
		tgtCol, tgtPiece := b.GetPieceAt(m.To)
		if tgtCol == None || tgtPiece == Empty {
			return fmt.Errorf("capture move with no piece at target square: %v", m.To)
		}
		// Remove piece currently on target square and move piece to taht position
		b.Pieces[tgtCol][tgtPiece] = b.Pieces[tgtCol][tgtPiece].Clear(m.To)
		b.movePiece(pCol, piece, m.From, m.To)
	case Promotion:
		b.Pieces[pCol][Pawn] = b.Pieces[pCol][Pawn].Clear(m.From)
		b.Pieces[pCol][m.Promotion] = b.Pieces[pCol][m.Promotion].Set(m.To)
		// TODO: Cover EnPassant, castling and Promotion + Capture cases
	default:
		return fmt.Errorf("unsupported move type: %v", m.Type)
	}

	b.UpdateOccupiedSquares()
	if b.SideToMove == Black {
		b.FullMoveCount += 1
	}

	b.SideToMove ^= 1 // toggle active player

	return nil
}

// PlayMoveSequence plays a sequence of moves, assuming alternating turns
func (b *Board) PlayMoveSequence(ml []Move) error {
	for _, m := range ml {
		err := b.PlayMove(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Board) isEqualBoard(other Board) bool {
	if b.OccupiedSquares != other.OccupiedSquares {
		return false
	}
	for color := White; color <= Black; color++ {
		for piece := Pawn; piece <= King; piece++ {
			if b.Pieces[color][piece] != other.Pieces[color][piece] {
				return false
			}
		}
	}
	return true
}

func (b *Board) movePiece(cl Color, p Piece, from, to Square) {
	b.Pieces[cl][p] = b.Pieces[cl][p].Clear(from).Set(to)
}

func (b *Board) ToFEN() string {
	enPassTgt := "-" // tracks en passant target square
	castAb := "KQkq" // tracks castling privileges
	hmClock := 0     // tracks the half move related to the 50 move draw rule
	sideToMove := "w"
	if b.SideToMove == Black {
		sideToMove = "b"
	}

	var sb strings.Builder
	// Board state, traverve ranks in reverse (8 to 1)
	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0
		// Traverse files 1 to 8
		for file := 0; file < 8; file++ {
			sq := Square(rank*8 + file)
			color, piece := b.GetPieceAt(sq)
			if piece == Empty {
				emptyCount++
				// If we are at the end of a rank, append the count
				if file == 7 && emptyCount > 0 {
					sb.WriteString(strconv.Itoa(emptyCount))
				}
				continue
			}
			// If we had empty squares before this piece, add the count
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
				emptyCount = 0
			}
			pieceStr := piece.String()
			if color == Black {
				pieceStr = strings.ToLower(pieceStr)
			}
			sb.WriteString(pieceStr)
		}
		if rank > 0 {
			sb.WriteRune('/')
		}
	}
	return fmt.Sprintf("%s %s %s %s %d %d", sb.String(), sideToMove, castAb, enPassTgt, hmClock, b.FullMoveCount)
}
