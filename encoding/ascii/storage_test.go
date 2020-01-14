package ascii

import (
	"reflect"
	"testing"

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
	encoder := NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		margins,
		models.White,
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

	if encoder.topColor != models.White {
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
		topColor    models.Color
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
				topColor:    models.Black,
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
				topColor:    models.White,
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
				topColor: models.Black,
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
				topColor: models.Black,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "\n" +
				"8rxxxkxxr\n" +
				"\n" +
				"\n" +
				"\n" +
				"7pxppqpbx\n" +
				"\n" +
				"\n" +
				"\n" +
				"6bnxxpnpx\n" +
				"\n" +
				"\n" +
				"\n" +
				"5xxxPNxxx\n" +
				"\n" +
				"\n" +
				"\n" +
				"4xpxxPxxx\n" +
				"\n" +
				"\n" +
				"\n" +
				"3xxNxxQxp\n" +
				"\n" +
				"\n" +
				"\n" +
				"2PPPBBPPP\n" +
				"\n" +
				"\n" +
				"\n" +
				"1RxxxKxxR\n" +
				"\n" +
				"\n" +
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
				topColor: models.White,
			},
			args: args{
				boardInFEN: kiwipete,
			},
			want: "\n" +
				"1RxxxKxxR\n" +
				"\n" +
				"\n" +
				"\n" +
				"2PPPBBPPP\n" +
				"\n" +
				"\n" +
				"\n" +
				"3xxNxxQxp\n" +
				"\n" +
				"\n" +
				"\n" +
				"4xpxxPxxx\n" +
				"\n" +
				"\n" +
				"\n" +
				"5xxxPNxxx\n" +
				"\n" +
				"\n" +
				"\n" +
				"6bnxxpnpx\n" +
				"\n" +
				"\n" +
				"\n" +
				"7pxppqpbx\n" +
				"\n" +
				"\n" +
				"\n" +
				"8rxxxkxxr\n" +
				"\n" +
				"\n" +
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
				topColor: models.Black,
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
			topColor:    data.fields.topColor,
		}
		got :=
			encoder.EncodePieceStorage(storage)

		if got != data.want {
			test.Log(got)
			test.Log(data.want)
			test.Fail()
		}
	}
}
