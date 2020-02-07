package ascii

// HorizontalMargins ...
type HorizontalMargins struct {
	Left  int
	Right int
}

// Width ...
func (margins HorizontalMargins) Width(
	contentWidth int,
) int {
	return margins.Left +
		margins.Right +
		contentWidth
}

// VerticalMargins ...
type VerticalMargins struct {
	Top    int
	Bottom int
}

// PieceMargins ...
type PieceMargins struct {
	HorizontalMargins
	VerticalMargins
}

// LegendMargins ...
type LegendMargins struct {
	File VerticalMargins
	Rank HorizontalMargins
}

// Margins ...
type Margins struct {
	Piece  PieceMargins
	Legend LegendMargins
}
