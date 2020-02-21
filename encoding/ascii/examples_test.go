package ascii_test

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func ExampleDecodeColor() {
	color, _ := ascii.DecodeColor("white")
	fmt.Printf("%v\n", color)

	// Output: 1
}

func ExampleEncodeColor() {
	color := ascii.EncodeColor(models.White)
	fmt.Printf("%v\n", color)

	// Output: white
}

func ExamplePieceStorageEncoder_EncodePieceStorage() {
	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		ascii.Margins{},
		ascii.WithoutColor,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// 8rxxxkxxr
	// 7pxppqpbx
	// 6bnxxpnpx
	// 5xxxPNxxx
	// 4xpxxPxxx
	// 3xxNxxQxp
	// 2PPPBBPPP
	// 1RxxxKxxR
	//  abcdefgh
}

func ExamplePieceStorageEncoder_EncodePieceStorage_withMargins() {
	margins := ascii.Margins{
		Piece: ascii.PieceMargins{
			HorizontalMargins: ascii.HorizontalMargins{
				Left: 1,
			},
		},
		Legend: ascii.LegendMargins{
			Rank: ascii.HorizontalMargins{
				Right: 1,
			},
		},
	}

	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		margins,
		ascii.WithoutColor,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// 8  r x x x k x x r
	// 7  p x p p q p b x
	// 6  b n x x p n p x
	// 5  x x x P N x x x
	// 4  x p x x P x x x
	// 3  x x N x x Q x p
	// 2  P P P B B P P P
	// 1  R x x x K x x R
	//    a b c d e f g h
}

func ExamplePieceStorageEncoder_EncodePieceStorage_withColors() {
	colorizer := func(text string, color climodels.OptionalColor) string {
		var colorMark byte
		if color.IsSet {
			colorMark = ascii.EncodeColor(color.Value)[0]
		} else {
			colorMark = 'n'
		}

		return fmt.Sprintf("(%c%s)", colorMark, text)
	}

	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		ascii.Margins{},
		colorizer,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// (n8)(wr)(bx)(wx)(bx)(wk)(bx)(wx)(br)
	// (n7)(bp)(wx)(bp)(wp)(bq)(wp)(bb)(wx)
	// (n6)(wb)(bn)(wx)(bx)(wp)(bn)(wp)(bx)
	// (n5)(bx)(wx)(bx)(wP)(bN)(wx)(bx)(wx)
	// (n4)(wx)(bp)(wx)(bx)(wP)(bx)(wx)(bx)
	// (n3)(bx)(wx)(bN)(wx)(bx)(wQ)(bx)(wp)
	// (n2)(wP)(bP)(wP)(bB)(wB)(bP)(wP)(bP)
	// (n1)(bR)(wx)(bx)(wx)(bK)(wx)(bx)(wR)
	// (n )(na)(nb)(nc)(nd)(ne)(nf)(ng)(nh)
}
