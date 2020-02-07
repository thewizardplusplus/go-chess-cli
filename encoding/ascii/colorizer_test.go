package ascii

import (
	"testing"

	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
)

func TestWithoutColor(test *testing.T) {
	type args struct {
		text  string
		color climodels.OptionalColor
	}
	type data struct {
		args args
		want string
	}

	for _, data := range []data{
		data{
			args: args{
				text: "test",
				color: climodels.NewOptionalColor(
					models.Black,
				),
			},
			want: "test",
		},
		data{
			args: args{
				text: "test",
				color: climodels.NewOptionalColor(
					models.White,
				),
			},
			want: "test",
		},
		data{
			args: args{
				text:  "test",
				color: climodels.WithoutColor,
			},
			want: "test",
		},
	} {
		got := WithoutColor(
			data.args.text,
			data.args.color,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
