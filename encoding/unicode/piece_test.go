package unicode

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

type MockPiece struct {
	kind  models.Kind
	color models.Color
}

func (piece MockPiece) Kind() models.Kind {
	return piece.kind
}

func (piece MockPiece) Color() models.Color {
	return piece.color
}

func (
	piece MockPiece,
) Position() models.Position {
	panic("not implemented")
}

func (piece MockPiece) ApplyPosition(
	position models.Position,
) models.Piece {
	panic("not implemented")
}

func (piece MockPiece) CheckMove(
	move models.Move,
	storage models.PieceStorage,
) bool {
	panic("not implemented")
}

func (piece MockPiece) String() string {
	panic("not implemented")
}

func TestEncodePiece(test *testing.T) {
	type args struct {
		piece models.Piece
	}
	type data struct {
		args args
		want string
	}

	for _, data := range []data{
		data{
			args: args{
				piece: MockPiece{
					kind:  models.King,
					color: models.White,
				},
			},
			want: "\u2654",
		},
		data{
			args: args{
				piece: MockPiece{
					kind:  models.Queen,
					color: models.Black,
				},
			},
			want: "\u265B",
		},
	} {
		got := EncodePiece(data.args.piece)

		if got != data.want {
			test.Fail()
		}
	}
}
