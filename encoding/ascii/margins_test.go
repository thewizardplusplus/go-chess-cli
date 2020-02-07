package ascii

import (
	"testing"
)

func TestHorizontalMarginsWidth(
	test *testing.T,
) {
	type fields struct {
		left  int
		right int
	}
	type args struct {
		contentWidth int
	}
	type data struct {
		fields fields
		args   args
		want   int
	}

	for _, data := range []data{
		data{
			fields: fields{2, 3},
			args:   args{0},
			want:   5,
		},
		data{
			fields: fields{2, 3},
			args:   args{4},
			want:   9,
		},
	} {
		margins := HorizontalMargins{
			Left:  data.fields.left,
			Right: data.fields.right,
		}
		got :=
			margins.Width(data.args.contentWidth)

		if got != data.want {
			test.Fail()
		}
	}
}
