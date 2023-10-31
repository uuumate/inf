package value

import (
	"github.com/go-playground/assert"
	"testing"
)

type ClearSliceParam struct {
	name string
	args struct {
		slice *[]interface{}
	}
	want *[]interface{}
}

func Test_ClearSlice(t *testing.T) {
	slice1 := []interface{}{1, 2, 3}

	pp := []ClearSliceParam{
		{
			name: "1",
			args: struct{ slice *[]interface{} }{slice: &slice1},
			want: &[]interface{}{},
		},
	}

	for _, p := range pp {
		t.Run(p.name, func(t *testing.T) {
			ClearSlice(p.args.slice)
			assert.Equal(t, p.args.slice, *p.want)
		})
	}
}
