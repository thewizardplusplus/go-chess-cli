package ascii

import (
	climodels "github.com/thewizardplusplus/go-chess-cli/models"
)

// Colorizer ...
type Colorizer func(
	text string,
	color climodels.OptionalColor,
) string

// WithoutColor ...
func WithoutColor(
	text string,
	color climodels.OptionalColor,
) string {
	return text
}
