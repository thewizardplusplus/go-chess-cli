package models

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewOptionalColor(test *testing.T) {
	got := NewOptionalColor(models.White)

	if got.Color != models.White {
		test.Fail()
	}
	if !got.IsSet {
		test.Fail()
	}
}
