package ascii

import (
	"strconv"
	"strings"

	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// PieceEncoder ...
type PieceEncoder func(
	piece models.Piece,
) string

// Colorizer ...
type Colorizer func(
	text string,
	color climodels.OptionalColor,
) string

// PieceStorageEncoder ...
type PieceStorageEncoder struct {
	encoder     PieceEncoder
	placeholder string
	margins     Margins
	colorizer   Colorizer
	topColor    models.Color
}

// NewPieceStorageEncoder ...
func NewPieceStorageEncoder(
	encoder PieceEncoder,
	placeholder string,
	margins Margins,
	colorizer Colorizer,
	topColor models.Color,
) PieceStorageEncoder {
	return PieceStorageEncoder{
		encoder:     encoder,
		placeholder: placeholder,
		margins:     margins,
		colorizer:   colorizer,
		topColor:    topColor,
	}
}

// EncodePieceStorage ...
func (
	encoder PieceStorageEncoder,
) EncodePieceStorage(
	storage models.PieceStorage,
) string {
	pieceMargins := encoder.margins.Piece
	legendMargins := encoder.margins.Legend

	var ranks []string
	var currentRank string
	positions := storage.Size().Positions()
	startedColor := models.Black
	currentColor := startedColor
	for _, position := range positions {
		if len(currentRank) == 0 {
			currentRank += encoder.wrapWithSpaces(
				strconv.Itoa(position.Rank+1),
				legendMargins.Rank,
				climodels.WithoutColor,
			)
		}

		var encodedPiece string
		piece, ok := storage.Piece(position)
		if ok {
			encodedPiece = encoder.encoder(piece)
		} else {
			encodedPiece = encoder.placeholder
		}
		currentRank += encoder.wrapWithSpaces(
			encodedPiece,
			pieceMargins.HorizontalMargins,
			climodels.NewOptionalColor(
				currentColor,
			),
		)

		currentColor = currentColor.Negative()

		lastFile := storage.Size().Height - 1
		if position.File == lastFile {
			ranks = append(ranks, currentRank)
			currentRank = ""

			startedColor = startedColor.Negative()
			currentColor = startedColor
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
				pieceMargins.VerticalMargins,
			)...,
		)
	}

	legendRank := encoder.spaces(
		legendMargins.Rank.Left+
			legendMargins.Rank.Right+
			1,
		climodels.WithoutColor,
	)
	width := storage.Size().Width
	for i := 0; i < width; i++ {
		legendRank += encoder.wrapWithSpaces(
			string(i+97),
			pieceMargins.HorizontalMargins,
			climodels.WithoutColor,
		)
	}
	sparseRanks = append(
		sparseRanks,
		wrapWithEmptyLines(
			legendRank,
			legendMargins.File,
		)...,
	)

	return strings.Join(sparseRanks, "\n")
}

func (
	encoder PieceStorageEncoder,
) wrapWithSpaces(
	text string,
	margins HorizontalMargins,
	color climodels.OptionalColor,
) string {
	prefix :=
		encoder.spaces(margins.Left, color)
	suffix :=
		encoder.spaces(margins.Right, color)
	text = encoder.colorizer(text, color)
	return prefix + text + suffix
}

func (
	encoder PieceStorageEncoder,
) spaces(
	length int,
	color climodels.OptionalColor,
) string {
	text := strings.Repeat(" ", length)
	return encoder.colorizer(text, color)
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
