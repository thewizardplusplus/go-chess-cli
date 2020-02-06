package models

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// OptionalColor ...
type OptionalColor struct {
	Color models.Color
	IsSet bool
}

// WithoutColor ...
var WithoutColor OptionalColor

// NewOptionalColor ...
func NewOptionalColor(
	color models.Color,
) OptionalColor {
	return OptionalColor{color, true}
}

// Negative ...
func (
	color OptionalColor,
) Negative() OptionalColor {
	if !color.IsSet {
		return WithoutColor
	}

	negativeColor := color.Color.Negative()
	return NewOptionalColor(negativeColor)
}
