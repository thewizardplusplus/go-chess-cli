package ascii

import (
	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

// Colorizer ...
type Colorizer func(
	text string,
	color models.Color,
) string

// OptionalColorizer ...
type OptionalColorizer func(
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

// NewOptionalColorizer ...
func NewOptionalColorizer(
	colorizer Colorizer,
) OptionalColorizer {
	return func(
		text string,
		color climodels.OptionalColor,
	) string {
		if !color.IsSet {
			return text
		}

		return colorizer(text, color.Value)
	}
}
