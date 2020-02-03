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
