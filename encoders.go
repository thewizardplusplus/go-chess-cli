package chesscli

import (
	"strconv"
	"strings"

	models "github.com/thewizardplusplus/go-chess-models"
)

// PieceEncoder ...
type PieceEncoder func(
	piece models.Piece,
) string

// PieceStorageEncoder ...
type PieceStorageEncoder struct {
	PieceEncoder PieceEncoder
	Separator    string
	TopColor     models.Color
}

// Encode ...
func (encoder PieceStorageEncoder) Encode(
	storage models.PieceStorage,
) string {
	var ranks []string
	var currentRank string
	positions := storage.Size().Positions()
	for _, position := range positions {
		if len(currentRank) == 0 {
			currentRank +=
				strconv.Itoa(position.Rank + 1)
		}

		piece, ok := storage.Piece(position)
		if ok {
			currentRank +=
				encoder.PieceEncoder(piece)
		} else {
			currentRank += encoder.Separator
		}

		lastFile := storage.Size().Height - 1
		if position.File == lastFile {
			ranks = append(ranks, currentRank)
			currentRank = ""
		}
	}
	if encoder.TopColor == models.Black {
		reverse(ranks)
	}

	legendRank := " "
	width := storage.Size().Width
	for i := 0; i < width; i++ {
		legendRank += string(i + 97)
	}
	ranks = append(ranks, legendRank)

	return strings.Join(ranks, "\n")
}

func reverse(strings []string) {
	left, right := 0, len(strings)-1
	for left < right {
		strings[left], strings[right] =
			strings[right], strings[left]
		left, right = left+1, right-1
	}
}
