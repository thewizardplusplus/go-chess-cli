package models

import (
	"reflect"
	"testing"

	models "github.com/thewizardplusplus/go-chess-models"
)

func TestNewOptionalColor(test *testing.T) {
	got := NewOptionalColor(models.White)

	if got.Value != models.White {
		test.Fail()
	}
	if !got.IsSet {
		test.Fail()
	}
}

func TestOptionalColorNegative(test *testing.T) {
	type fields struct {
		value models.Color
		isSet bool
	}
	type data struct {
		fields fields
		want   OptionalColor
	}

	for _, data := range []data{
		{
			fields: fields{models.Black, true},
			want: OptionalColor{
				Value: models.White,
				IsSet: true,
			},
		},
		{
			fields: fields{models.White, true},
			want: OptionalColor{
				Value: models.Black,
				IsSet: true,
			},
		},
		{
			fields: fields{isSet: false},
			want:   WithoutColor,
		},
	} {
		color := OptionalColor{
			Value: data.fields.value,
			IsSet: data.fields.isSet,
		}
		got := color.Negative()

		if !reflect.DeepEqual(got, data.want) {
			test.Fail()
		}
	}
}
