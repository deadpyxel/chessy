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
		// Clear the piece from its original position
		b.Pieces[pCol][piece] = b.Pieces[pCol][piece].Clear(m.From)
		// Sets the piece in its new position
		b.Pieces[pCol][piece] = b.Pieces[pCol][piece].Set(m.To)
	}

	b.UpdateOccupiedSquares()

	b.SideToMove ^= 1 // toggle active player

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
