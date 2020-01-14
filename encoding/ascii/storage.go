package ascii

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
	encoder     PieceEncoder
	placeholder string
	margins     Margins
	topColor    models.Color
}

// NewPieceStorageEncoder ...
func NewPieceStorageEncoder(
	encoder PieceEncoder,
	placeholder string,
	margins Margins,
	topColor models.Color,
) PieceStorageEncoder {
	return PieceStorageEncoder{
		encoder:     encoder,
		placeholder: placeholder,
		margins:     margins,
		topColor:    topColor,
	}
}

// EncodePieceStorage ...
func (
	encoder PieceStorageEncoder,
) EncodePieceStorage(
	storage models.PieceStorage,
) string {
	var ranks []string
	var currentRank string
	positions := storage.Size().Positions()
	for _, position := range positions {
		if len(currentRank) == 0 {
			currentRank += spaces(
				encoder.margins.Legend.Rank.Left,
			)
			currentRank +=
				strconv.Itoa(position.Rank + 1)
			currentRank += spaces(
				encoder.margins.Legend.Rank.Right,
			)
		}

		currentRank +=
			spaces(encoder.margins.Piece.Left)

		piece, ok := storage.Piece(position)
		if ok {
			currentRank += encoder.encoder(piece)
		} else {
			currentRank += encoder.placeholder
		}

		currentRank +=
			spaces(encoder.margins.Piece.Right)

		lastFile := storage.Size().Height - 1
		if position.File == lastFile {
			ranks = append(ranks, currentRank)
			currentRank = ""
		}
	}
	if encoder.topColor == models.Black {
		reverse(ranks)
	}

	var sparseRanks []string
	for _, rank := range ranks {
		sparseRanks = append(
			sparseRanks,
			empties(encoder.margins.Piece.Top)...,
		)
		sparseRanks = append(sparseRanks, rank)
		sparseRanks = append(
			sparseRanks,
			empties(
				encoder.margins.Piece.Bottom,
			)...,
		)
	}

	legendRank := spaces(
		encoder.margins.Legend.Rank.Left +
			encoder.margins.Legend.Rank.Right +
			1,
	)
	width := storage.Size().Width
	for i := 0; i < width; i++ {
		legendRank +=
			spaces(encoder.margins.Piece.Left) +
				string(i+97) +
				spaces(encoder.margins.Piece.Right)
	}
	sparseRanks =
		append(sparseRanks, legendRank)

	return strings.Join(sparseRanks, "\n")
}

func spaces(length int) string {
	return strings.Repeat(" ", length)
}

func empties(count int) []string {
	return make([]string, count)
}

func reverse(strings []string) {
	left, right := 0, len(strings)-1
	for left < right {
		strings[left], strings[right] =
			strings[right], strings[left]
		left, right = left+1, right-1
	}
}
