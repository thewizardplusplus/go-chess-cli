package chesscli

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

const (
	kiwipete = "r3k2r/p1ppqpb1/bn2pnp1/3PN3" +
		"/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
)

func TestPieceStorageEncoderEncode(
	test *testing.T,
) {
	type fields struct {
		pieceEncoder PieceEncoder
		separator    string
		topColor     models.Color
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
				pieceEncoder: uci.EncodePiece,
				separator:    "x",
				topColor:     models.Black,
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
				pieceEncoder: uci.EncodePiece,
				separator:    "x",
				topColor:     models.White,
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
			PieceEncoder: data.fields.pieceEncoder,
			Separator:    data.fields.separator,
			TopColor:     data.fields.topColor,
		}
		got := encoder.Encode(storage)

		if got != data.want {
			test.Fail()
		}
	}
}
