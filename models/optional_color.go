package models

import (
	models "github.com/thewizardplusplus/go-chess-models"
)

// OptionalColor ...
type OptionalColor struct {
	Value models.Color
	IsSet bool
}

// ...
var (
	WithoutColor OptionalColor
)

// NewOptionalColor ...
func NewOptionalColor(color models.Color) OptionalColor {
	return OptionalColor{color, true}
}

// Negative ...
func (color OptionalColor) Negative() OptionalColor {
	if !color.IsSet {
		return WithoutColor
	}

	negativeColor := color.Value.Negative()
	return NewOptionalColor(negativeColor)
}
