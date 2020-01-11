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
		data{
			args:      args{"black"},
			wantColor: models.Black,
			wantErr:   false,
		},
		data{
			args:      args{"white"},
			wantColor: models.White,
			wantErr:   false,
		},
		data{
			args:      args{"incorrect"},
			wantColor: 0,
			wantErr:   true,
		},
	} {
		gotColor, gotErr :=
			DecodeColor(data.args.text)

		if gotColor != data.wantColor {
			test.Fail()
		}
		hasErr := gotErr != nil
		if hasErr != data.wantErr {
			test.Fail()
		}
	}
}
