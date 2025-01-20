package board

import "fmt"

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
