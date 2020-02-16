package ascii

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

const (
	kiwipete = "r3k2r/p1ppqpb1/bn2pnp1/3PN3" +
		"/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
)

func TestNewPieceStorageEncoder(
	test *testing.T,
) {
	margins := Margins{
		Piece: PieceMargins{
			HorizontalMargins: HorizontalMargins{
				Left:  1,
				Right: 2,
			},
			VerticalMargins: VerticalMargins{
				Top:    3,
				Bottom: 4,
			},
		},
		Legend: LegendMargins{
			File: VerticalMargins{
				Top:    5,
				Bottom: 6,
			},
			Rank: HorizontalMargins{
				Left:  7,
				Right: 8,
			},
		},
	}
	colorizer := func(
		text string,
		color climodels.OptionalColor,
	) string {
		panic("not implemented")
	}
	encoder := NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		margins,
		colorizer,
		models.White,
		2,
	)

	gotEncoder := reflect.
		ValueOf(encoder.encoder).
		Pointer()
	wantEncoder := reflect.
		ValueOf(uci.EncodePiece).
		Pointer()
	if gotEncoder != wantEncoder {
		test.Fail()
	}

	if encoder.placeholder != "x" {
		test.Fail()
	}

	if !reflect.DeepEqual(
		encoder.margins,
		margins,
	) {
		test.Fail()
	}

	gotColorizer := reflect.
		ValueOf(encoder.colorizer).
		Pointer()
	wantColorizer := reflect.
		ValueOf(colorizer).
		Pointer()
	if gotColorizer != wantColorizer {
		test.Fail()
	}

	if encoder.topColor != models.White {
		test.Fail()
	}

	if encoder.pieceWidth != 2 {
		test.Fail()
	}
}

func TestPieceStorageEncoderEncodePieceStorage(
	test *testing.T,
) {
	type fields struct {
		encoder     PieceEncoder
		placeholder string
		margins     Margins
		colorizer   OptionalColorizer
		topColor    models.Color
		pieceWidth  int
	}
	type args struct {
		boardInFEN string
	}
	type data struct {
		fields fields
		args   args
		want   string
	}

	for _, data := range []data{
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins:     Margins{},
				colorizer:   WithoutColor,
				topColor:    models.Black,
				pieceWidth:  1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "8rxxxkxxr\n" +
				"7pxppqpbx\n" +
				"6bnxxpnpx\n" +
				"5xxxPNxxx\n" +
				"4xpxxPxxx\n" +
				"3xxNxxQxp\n" +
				"2PPPBBPPP\n" +
				"1RxxxKxxR\n" +
				" abcdefgh",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins:     Margins{},
				colorizer:   WithoutColor,
				topColor:    models.White,
				pieceWidth:  1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "1RxxxKxxR\n" +
				"2PPPBBPPP\n" +
				"3xxNxxQxp\n" +
				"4xpxxPxxx\n" +
				"5xxxPNxxx\n" +
				"6bnxxpnpx\n" +
				"7pxppqpbx\n" +
				"8rxxxkxxr\n" +
				" abcdefgh",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "8 r   x   x   x   k   x   x   r  \n" +
				"7 p   x   p   p   q   p   b   x  \n" +
				"6 b   n   x   x   p   n   p   x  \n" +
				"5 x   x   x   P   N   x   x   x  \n" +
				"4 x   p   x   x   P   x   x   x  \n" +
				"3 x   x   N   x   x   Q   x   p  \n" +
				"2 P   P   P   B   B   P   P   P  \n" +
				"1 R   x   x   x   K   x   x   R  \n" +
				"  a   b   c   d   e   f   g   h  ",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: strings.Repeat(" ", 9) + "\n" +
				"8rxxxkxxr\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"7pxppqpbx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"6bnxxpnpx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"5xxxPNxxx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"4xpxxPxxx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"3xxNxxQxp\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"2PPPBBPPP\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"1RxxxKxxR\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				" abcdefgh",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.White,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: strings.Repeat(" ", 9) + "\n" +
				"1RxxxKxxR\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"2PPPBBPPP\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"3xxNxxQxp\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"4xpxxPxxx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"5xxxPNxxx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"6bnxxpnpx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"7pxppqpbx\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				"8rxxxkxxr\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9) + "\n" +
				" abcdefgh",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Legend: LegendMargins{
						Rank: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: " 8  rxxxkxxr\n" +
				" 7  pxppqpbx\n" +
				" 6  bnxxpnpx\n" +
				" 5  xxxPNxxx\n" +
				" 4  xpxxPxxx\n" +
				" 3  xxNxxQxp\n" +
				" 2  PPPBBPPP\n" +
				" 1  RxxxKxxR\n" +
				"    abcdefgh",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Legend: LegendMargins{
						File: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "8rxxxkxxr\n" +
				"7pxppqpbx\n" +
				"6bnxxpnpx\n" +
				"5xxxPNxxx\n" +
				"4xpxxPxxx\n" +
				"3xxNxxQxp\n" +
				"2PPPBBPPP\n" +
				"1RxxxKxxR\n" +
				strings.Repeat(" ", 9) + "\n" +
				" abcdefgh\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9),
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
					Legend: LegendMargins{
						File: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
						Rank: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: strings.Repeat(" ", 4*9) + "\n" +
				" 8   r   x   x   x   k   x   x   r  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 7   p   x   p   p   q   p   b   x  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 6   b   n   x   x   p   n   p   x  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 5   x   x   x   P   N   x   x   x  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 4   x   p   x   x   P   x   x   x  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 3   x   x   N   x   x   Q   x   p  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 2   P   P   P   B   B   P   P   P  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				" 1   R   x   x   x   K   x   x   R  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9) + "\n" +
				"     a   b   c   d   e   f   g   h  \n" +
				strings.Repeat(" ", 4*9) + "\n" +
				strings.Repeat(" ", 4*9),
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins:     Margins{},
				colorizer: func(
					text string,
					color climodels.OptionalColor,
				) string {
					var colorMark byte
					if color.IsSet {
						colorMark =
							EncodeColor(color.Value)[0]
					} else {
						colorMark = 'n'
					}

					return fmt.Sprintf(
						"(%c%s)",
						colorMark,
						text,
					)
				},
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "(n8)(wr)(bx)(wx)(bx)(wk)(bx)(wx)(br)\n" +
				"(n7)(bp)(wx)(bp)(wp)(bq)(wp)(bb)(wx)\n" +
				"(n6)(wb)(bn)(wx)(bx)(wp)(bn)(wp)(bx)\n" +
				"(n5)(bx)(wx)(bx)(wP)(bN)(wx)(bx)(wx)\n" +
				"(n4)(wx)(bp)(wx)(bx)(wP)(bx)(wx)(bx)\n" +
				"(n3)(bx)(wx)(bN)(wx)(bx)(wQ)(bx)(wp)\n" +
				"(n2)(wP)(bP)(wP)(bB)(wB)(bP)(wP)(bP)\n" +
				"(n1)(bR)(wx)(bx)(wx)(bK)(wx)(bx)(wR)\n" +
				"(n )(na)(nb)(nc)(nd)(ne)(nf)(ng)(nh)",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 1,
						},
					},
					Legend: LegendMargins{
						Rank: HorizontalMargins{
							Left:  1,
							Right: 1,
						},
					},
				},
				colorizer: func(
					text string,
					color climodels.OptionalColor,
				) string {
					var colorMark byte
					if color.IsSet {
						colorMark =
							EncodeColor(color.Value)[0]
					} else {
						colorMark = 'n'
					}

					return fmt.Sprintf(
						"(%c%s)",
						colorMark,
						text,
					)
				},
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "(n )(n8)(n )(w )(wr)(w )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )(w )(wk)(w )(b )(bx)(b )(w )(wx)(w )(b )(br)(b )\n" +
				"(n )(n7)(n )(b )(bp)(b )(w )(wx)(w )(b )(bp)(b )(w )(wp)(w )(b )(bq)(b )(w )(wp)(w )(b )(bb)(b )(w )(wx)(w )\n" +
				"(n )(n6)(n )(w )(wb)(w )(b )(bn)(b )(w )(wx)(w )(b )(bx)(b )(w )(wp)(w )(b )(bn)(b )(w )(wp)(w )(b )(bx)(b )\n" +
				"(n )(n5)(n )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )(w )(wP)(w )(b )(bN)(b )(w )(wx)(w )(b )(bx)(b )(w )(wx)(w )\n" +
				"(n )(n4)(n )(w )(wx)(w )(b )(bp)(b )(w )(wx)(w )(b )(bx)(b )(w )(wP)(w )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )\n" +
				"(n )(n3)(n )(b )(bx)(b )(w )(wx)(w )(b )(bN)(b )(w )(wx)(w )(b )(bx)(b )(w )(wQ)(w )(b )(bx)(b )(w )(wp)(w )\n" +
				"(n )(n2)(n )(w )(wP)(w )(b )(bP)(b )(w )(wP)(w )(b )(bB)(b )(w )(wB)(w )(b )(bP)(b )(w )(wP)(w )(b )(bP)(b )\n" +
				"(n )(n1)(n )(b )(bR)(b )(w )(wx)(w )(b )(bx)(b )(w )(wx)(w )(b )(bK)(b )(w )(wx)(w )(b )(bx)(b )(w )(wR)(w )\n" +
				"(n   )(n )(na)(n )(n )(nb)(n )(n )(nc)(n )(n )(nd)(n )(n )(ne)(n )(n )(nf)(n )(n )(ng)(n )(n )(nh)(n )",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 1,
						},
					},
					Legend: LegendMargins{
						File: VerticalMargins{
							Top:    1,
							Bottom: 1,
						},
					},
				},
				colorizer: func(
					text string,
					color climodels.OptionalColor,
				) string {
					var colorMark byte
					if color.IsSet {
						colorMark =
							EncodeColor(color.Value)[0]
					} else {
						colorMark = 'n'
					}

					return fmt.Sprintf(
						"(%c%s)",
						colorMark,
						text,
					)
				},
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n8)(wr)(bx)(wx)(bx)(wk)(bx)(wx)(br)\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n7)(bp)(wx)(bp)(wp)(bq)(wp)(bb)(wx)\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n6)(wb)(bn)(wx)(bx)(wp)(bn)(wp)(bx)\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n5)(bx)(wx)(bx)(wP)(bN)(wx)(bx)(wx)\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n4)(wx)(bp)(wx)(bx)(wP)(bx)(wx)(bx)\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n3)(bx)(wx)(bN)(wx)(bx)(wQ)(bx)(wp)\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n2)(wP)(bP)(wP)(bB)(wB)(bP)(wP)(bP)\n" +
				"(n )(w )(b )(w )(b )(w )(b )(w )(b )\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n1)(bR)(wx)(bx)(wx)(bK)(wx)(bx)(wR)\n" +
				"(n )(b )(w )(b )(w )(b )(w )(b )(w )\n" +
				"(n )(n )(n )(n )(n )(n )(n )(n )(n )\n" +
				"(n )(na)(nb)(nc)(nd)(ne)(nf)(ng)(nh)\n" +
				"(n )(n )(n )(n )(n )(n )(n )(n )(n )",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Piece: PieceMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 1,
						},
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 1,
						},
					},
					Legend: LegendMargins{
						File: VerticalMargins{
							Top:    1,
							Bottom: 1,
						},
						Rank: HorizontalMargins{
							Left:  1,
							Right: 1,
						},
					},
				},
				colorizer: func(
					text string,
					color climodels.OptionalColor,
				) string {
					var colorMark byte
					if color.IsSet {
						colorMark =
							EncodeColor(color.Value)[0]
					} else {
						colorMark = 'n'
					}

					return fmt.Sprintf(
						"(%c%s)",
						colorMark,
						text,
					)
				},
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n )(n8)(n )(w )(wr)(w )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )(w )(wk)(w )(b )(bx)(b )(w )(wx)(w )(b )(br)(b )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n )(n7)(n )(b )(bp)(b )(w )(wx)(w )(b )(bp)(b )(w )(wp)(w )(b )(bq)(b )(w )(wp)(w )(b )(bb)(b )(w )(wx)(w )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n )(n6)(n )(w )(wb)(w )(b )(bn)(b )(w )(wx)(w )(b )(bx)(b )(w )(wp)(w )(b )(bn)(b )(w )(wp)(w )(b )(bx)(b )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n )(n5)(n )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )(w )(wP)(w )(b )(bN)(b )(w )(wx)(w )(b )(bx)(b )(w )(wx)(w )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n )(n4)(n )(w )(wx)(w )(b )(bp)(b )(w )(wx)(w )(b )(bx)(b )(w )(wP)(w )(b )(bx)(b )(w )(wx)(w )(b )(bx)(b )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n )(n3)(n )(b )(bx)(b )(w )(wx)(w )(b )(bN)(b )(w )(wx)(w )(b )(bx)(b )(w )(wQ)(w )(b )(bx)(b )(w )(wp)(w )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n )(n2)(n )(w )(wP)(w )(b )(bP)(b )(w )(wP)(w )(b )(bB)(b )(w )(wB)(w )(b )(bP)(b )(w )(wP)(w )(b )(bP)(b )\n" +
				"(n   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n )(n1)(n )(b )(bR)(b )(w )(wx)(w )(b )(bx)(b )(w )(wx)(w )(b )(bK)(b )(w )(wx)(w )(b )(bx)(b )(w )(wR)(w )\n" +
				"(n   )(b   )(w   )(b   )(w   )(b   )(w   )(b   )(w   )\n" +
				"(n   )(n   )(n   )(n   )(n   )(n   )(n   )(n   )(n   )\n" +
				"(n   )(n )(na)(n )(n )(nb)(n )(n )(nc)(n )(n )(nd)(n )(n )(ne)(n )(n )(nf)(n )(n )(ng)(n )(n )(nh)(n )\n" +
				"(n   )(n   )(n   )(n   )(n   )(n   )(n   )(n   )(n   )",
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Board: VerticalMargins{
						Top:    1,
						Bottom: 2,
					},
				},
				colorizer:  WithoutColor,
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: strings.Repeat(" ", 9) + "\n" +
				"8rxxxkxxr\n" +
				"7pxppqpbx\n" +
				"6bnxxpnpx\n" +
				"5xxxPNxxx\n" +
				"4xpxxPxxx\n" +
				"3xxNxxQxp\n" +
				"2PPPBBPPP\n" +
				"1RxxxKxxR\n" +
				" abcdefgh\n" +
				strings.Repeat(" ", 9) + "\n" +
				strings.Repeat(" ", 9),
		},
		data{
			fields: fields{
				encoder:     uci.EncodePiece,
				placeholder: "x",
				margins: Margins{
					Board: VerticalMargins{
						Top:    1,
						Bottom: 2,
					},
				},
				colorizer: func(
					text string,
					color climodels.OptionalColor,
				) string {
					var colorMark byte
					if color.IsSet {
						colorMark =
							EncodeColor(color.Value)[0]
					} else {
						colorMark = 'n'
					}

					return fmt.Sprintf(
						"(%c%s)",
						colorMark,
						text,
					)
				},
				topColor:   models.Black,
				pieceWidth: 1,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "(n )(n )(n )(n )(n )(n )(n )(n )(n )\n" +
				"(n8)(wr)(bx)(wx)(bx)(wk)(bx)(wx)(br)\n" +
				"(n7)(bp)(wx)(bp)(wp)(bq)(wp)(bb)(wx)\n" +
				"(n6)(wb)(bn)(wx)(bx)(wp)(bn)(wp)(bx)\n" +
				"(n5)(bx)(wx)(bx)(wP)(bN)(wx)(bx)(wx)\n" +
				"(n4)(wx)(bp)(wx)(bx)(wP)(bx)(wx)(bx)\n" +
				"(n3)(bx)(wx)(bN)(wx)(bx)(wQ)(bx)(wp)\n" +
				"(n2)(wP)(bP)(wP)(bB)(wB)(bP)(wP)(bP)\n" +
				"(n1)(bR)(wx)(bx)(wx)(bK)(wx)(bx)(wR)\n" +
				"(n )(na)(nb)(nc)(nd)(ne)(nf)(ng)(nh)\n" +
				"(n )(n )(n )(n )(n )(n )(n )(n )(n )\n" +
				"(n )(n )(n )(n )(n )(n )(n )(n )(n )",
		},
	} {
		storage, err := uci.DecodePieceStorage(
			data.args.boardInFEN,
			pieces.NewPiece,
			models.NewBoard,
		)
		if err != nil {
			test.Fail()
			continue
		}

		encoder := PieceStorageEncoder{
			encoder:     data.fields.encoder,
			placeholder: data.fields.placeholder,
			margins:     data.fields.margins,
			colorizer:   data.fields.colorizer,
			topColor:    data.fields.topColor,
			pieceWidth:  data.fields.pieceWidth,
		}
		got :=
			encoder.EncodePieceStorage(storage)

		if got != data.want {
			test.Fail()
		}
	}
}
