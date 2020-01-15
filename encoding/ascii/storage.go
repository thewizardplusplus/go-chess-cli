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
			currentRank += wrapWithSpaces(
				strconv.Itoa(position.Rank+1),
				encoder.margins.Legend.Rank,
			)
		}

		var encodedPiece string
		piece, ok := storage.Piece(position)
		if ok {
			encodedPiece = encoder.encoder(piece)
		} else {
			encodedPiece = encoder.placeholder
		}
		currentRank += wrapWithSpaces(
			encodedPiece,
			encoder.margins.
				Piece.HorizontalMargins,
		)

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
			wrapWithEmptyLines(
				rank,
				encoder.margins.
					Piece.VerticalMargins,
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
		legendRank += wrapWithSpaces(
			string(i+97),
			encoder.margins.
				Piece.HorizontalMargins,
		)
	}

	if encoder.margins.Legend.File.Bottom > 0 {
		encoder.margins.Legend.File.Bottom++
	}
	sparseRanks = append(
		sparseRanks,
		wrapWithEmptyLines(
			legendRank,
			encoder.margins.Legend.File,
		)...,
	)

	return strings.Join(sparseRanks, "\n")
}

func wrapWithSpaces(
	text string,
	margins HorizontalMargins,
) string {
	return spaces(margins.Left) +
		text +
		spaces(margins.Right)
}

func spaces(length int) string {
	return strings.Repeat(" ", length)
}

func wrapWithEmptyLines(
	line string,
	margins VerticalMargins,
) []string {
	var lines []string
	lines = append(
		lines,
		emptyLines(margins.Top)...,
	)
	lines = append(lines, line)
	lines = append(
		lines,
		emptyLines(margins.Bottom)...,
	)

	return lines
}

func emptyLines(count int) []string {
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
