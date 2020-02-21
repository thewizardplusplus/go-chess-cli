package ascii

import (
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestDecodeColor(test *testing.T) {
	type args struct {
		text string
	}
	type data struct {
		args      args
		wantColor models.Color
		wantErr   bool
	}

	for _, data := range []data{
		{
			args:      args{"black"},
			wantColor: models.Black,
			wantErr:   false,
		},
		{
			args:      args{"white"},
			wantColor: models.White,
			wantErr:   false,
		},
		{
			args:      args{"incorrect"},
			wantColor: 0,
			wantErr:   true,
		},
	} {
		gotColor, gotErr := DecodeColor(data.args.text)

		if gotColor != data.wantColor {
			test.Fail()
		}
		if hasErr := gotErr != nil; hasErr != data.wantErr {
			test.Fail()
		}
	}
}

func TestEncodeColor(test *testing.T) {
	type args struct {
		color models.Color
	}
	type data struct {
		args args
		want string
	}

	for _, data := range []data{
		{
			args: args{models.Black},
			want: "black",
		},
		{
			args: args{models.White},
			want: "white",
		},
	} {
		got := EncodeColor(data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}
