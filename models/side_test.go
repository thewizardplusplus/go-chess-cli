package models

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewSide(test *testing.T) {
	type args struct {
		color models.Color
	}
	type data struct {
		args args
		want Side
	}

	for _, data := range []data{
		data{
			args: args{models.Black},
			want: Searcher,
		},
		data{
			args: args{models.White},
			want: Human,
		},
	} {
		got := NewSide(data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestSideInvert(test *testing.T) {
	type data struct {
		side Side
		want Side
	}

	for _, data := range []data{
		data{
			side: Searcher,
			want: Human,
		},
		data{
			side: Human,
			want: Searcher,
		},
	} {
		got := data.side.Invert()

		if got != data.want {
			test.Fail()
		}
	}
}
