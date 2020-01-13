package unicode

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// EncodePiece ...
func EncodePiece(piece models.Piece) string {
	var text string
	switch piece.Kind() {
	case models.King:
		switch piece.Color() {
		case models.Black:
			text = "\u265a"
		case models.White:
			text = "\u2654"
		}
	case models.Queen:
		switch piece.Color() {
		case models.Black:
			text = "\u265b"
		case models.White:
			text = "\u2655"
		}
	case models.Rook:
		switch piece.Color() {
		case models.Black:
			text = "\u265c"
		case models.White:
			text = "\u2656"
		}
	case models.Bishop:
		switch piece.Color() {
		case models.Black:
			text = "\u265d"
		case models.White:
			text = "\u2657"
		}
	case models.Knight:
		switch piece.Color() {
		case models.Black:
			text = "\u265e"
		case models.White:
			text = "\u2658"
		}
	case models.Pawn:
		switch piece.Color() {
		case models.Black:
			text = "\u265f"
		case models.White:
			text = "\u2659"
		}
	}

	return text
}
