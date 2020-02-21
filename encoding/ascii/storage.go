package ascii

import (
	"strconv"
	"strings"

	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// PieceEncoder ...
type PieceEncoder func(piece models.Piece) string

// PieceStorageEncoder ...
type PieceStorageEncoder struct {
	encoder     PieceEncoder
	placeholder string
	margins     Margins
	colorizer   OptionalColorizer
	topColor    models.Color
	pieceWidth  int
}

// NewPieceStorageEncoder ...
func NewPieceStorageEncoder(
	encoder PieceEncoder,
	placeholder string,
	margins Margins,
	colorizer OptionalColorizer,
	topColor models.Color,
	pieceWidth int,
) PieceStorageEncoder {
	return PieceStorageEncoder{
		encoder:     encoder,
		placeholder: placeholder,
		margins:     margins,
		colorizer:   colorizer,
		topColor:    topColor,
		pieceWidth:  pieceWidth,
	}
}

// EncodePieceStorage ...
func (encoder PieceStorageEncoder) EncodePieceStorage(
	storage models.PieceStorage,
) string {
	pieceMargins := encoder.margins.Piece
	legendMargins := encoder.margins.Legend

	var ranks []string
	var currentRank string
	startedColor := models.Black
	currentColor := startedColor
	for _, position := range storage.Size().Positions() {
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
			climodels.NewOptionalColor(currentColor),
		)

		currentColor = currentColor.Negative()
		if lastFile := storage.Size().Height - 1; position.File == lastFile {
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
		startedColor = startedColor.Negative()
		sparseRanks = append(sparseRanks, encoder.wrapWithEmptyLines(
			[]string{rank},
			storage.Size().Width,
			pieceMargins.VerticalMargins,
			climodels.NewOptionalColor(startedColor),
		)...)
	}

	legendRank :=
		encoder.spaces(legendMargins.Rank.Width(1), climodels.WithoutColor)
	for i := 0; i < storage.Size().Width; i++ {
		legendRank += encoder.wrapWithSpaces(
			string(i+97),
			pieceMargins.HorizontalMargins,
			climodels.WithoutColor,
		)
	}
	sparseRanks = append(sparseRanks, encoder.wrapWithEmptyLines(
		[]string{legendRank},
		storage.Size().Width,
		legendMargins.File,
		climodels.WithoutColor,
	)...)

	sparseRanks = encoder.wrapWithEmptyLines(
		sparseRanks,
		storage.Size().Width,
		encoder.margins.Board,
		climodels.WithoutColor,
	)

	return strings.Join(sparseRanks, "\n")
}

func (encoder PieceStorageEncoder) wrapWithSpaces(
	text string,
	margins HorizontalMargins,
	color climodels.OptionalColor,
) string {
	prefix := encoder.spaces(margins.Left, color)
	suffix := encoder.spaces(margins.Right, color)
	text = encoder.colorizer(text, color)
	return prefix + text + suffix
}

func (encoder PieceStorageEncoder) spaces(
	length int,
	color climodels.OptionalColor,
) string {
	if length == 0 {
		return ""
	}

	text := strings.Repeat(" ", length)
	return encoder.colorizer(text, color)
}

func (encoder PieceStorageEncoder) wrapWithEmptyLines(
	lines []string,
	width int,
	margins VerticalMargins,
	startedColor climodels.OptionalColor,
) []string {
	var wrappedLines []string
	wrappedLines = append(wrappedLines, encoder.emptyLines(
		margins.Top,
		width,
		startedColor,
	)...)
	wrappedLines = append(wrappedLines, lines...)
	wrappedLines = append(wrappedLines, encoder.emptyLines(
		margins.Bottom,
		width,
		startedColor,
	)...)

	return wrappedLines
}

func (encoder PieceStorageEncoder) emptyLines(
	count int,
	width int,
	startedColor climodels.OptionalColor,
) []string {
	var lines []string
	for i := 0; i < count; i++ {
		line := encoder.emptyLine(width, startedColor)
		lines = append(lines, line)
	}

	return lines
}

func (encoder PieceStorageEncoder) emptyLine(
	width int,
	startedColor climodels.OptionalColor,
) string {
	pieceMargins := encoder.margins.Piece
	legendMargins := encoder.margins.Legend

	line := encoder.spaces(legendMargins.Rank.Width(1), climodels.WithoutColor)
	for i, currentColor := 0, startedColor; i < width; i++ {
		line += encoder.spaces(pieceMargins.Width(encoder.pieceWidth), currentColor)
		currentColor = currentColor.Negative()
	}

	return line
}

func reverse(strings []string) {
	left, right := 0, len(strings)-1
	for left < right {
		strings[left], strings[right] = strings[right], strings[left]
		left, right = left+1, right-1
	}
}
