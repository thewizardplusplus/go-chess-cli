package ascii

import (
	"fmt"
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
		{
			args: args{
				text:  "test",
				color: climodels.NewOptionalColor(models.Black),
			},
			want: "test",
		},
		{
			args: args{
				text:  "test",
				color: climodels.NewOptionalColor(models.White),
			},
			want: "test",
		},
		{
			args: args{
				text:  "test",
				color: climodels.WithoutColor,
			},
			want: "test",
		},
	} {
		got := WithoutColor(data.args.text, data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}

func TestNewOptionalColorizer(test *testing.T) {
	type fields struct {
		colorizer Colorizer
	}
	type args struct {
		text  string
		color climodels.OptionalColor
	}
	type data struct {
		fields fields
		args   args
		want   string
	}

	for _, data := range []data{
		{
			fields: fields{
				colorizer: func(text string, color models.Color) string {
					if text != "test" {
						test.Fail()
					}
					if color != models.Black {
						test.Fail()
					}

					return fmt.Sprintf("(%s:%s)", EncodeColor(color), text)
				},
			},
			args: args{
				text:  "test",
				color: climodels.NewOptionalColor(models.Black),
			},
			want: "(black:test)",
		},
		{
			fields: fields{
				colorizer: func(text string, color models.Color) string {
					if text != "test" {
						test.Fail()
					}
					if color != models.White {
						test.Fail()
					}

					return fmt.Sprintf("(%s:%s)", EncodeColor(color), text)
				},
			},
			args: args{
				text:  "test",
				color: climodels.NewOptionalColor(models.White),
			},
			want: "(white:test)",
		},
		{
			fields: fields{
				colorizer: func(text string, color models.Color) string {
					panic("not implemented")
				},
			},
			args: args{
				text:  "test",
				color: climodels.WithoutColor,
			},
			want: "test",
		},
	} {
		colorizer := NewOptionalColorizer(data.fields.colorizer)
		got := colorizer(data.args.text, data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}
