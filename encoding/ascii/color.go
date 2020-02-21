package ascii

import (
	"errors"

	models "github.com/thewizardplusplus/go-chess-models"
)

// DecodeColor ...
func DecodeColor(text string) (models.Color, error) {
	var color models.Color
	switch text {
	case "black":
		color = models.Black
	case "white":
		color = models.White
	default:
		return 0, errors.New("incorrect color")
	}

	return color, nil
}

// EncodeColor ...
func EncodeColor(color models.Color) string {
	var text string
	switch color {
	case models.Black:
		text = "black"
	case models.White:
		text = "white"
	}

	return text
}
