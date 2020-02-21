package unicode_test

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/unicode"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func ExampleEncodePiece() {
	fen := unicode.EncodePiece(pieces.NewBishop(models.White, models.Position{}))
	fmt.Printf("%v\n", fen)

	// Output: ♗
}
