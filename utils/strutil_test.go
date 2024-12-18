package utils_test

import (
	"reflect"
	"testing"

	"go.farcloser.world/core/utils"
)

func TestConvertKVStringsToMap(t *testing.T) {
	t.Parallel()

	type args struct {
		values []string
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "normal",
			args: args{
				values: []string{"foo=bar", "baz=qux"},
			},
			want: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
		},
		{
			name: "normal-1",
			args: args{
				values: []string{"foo"},
			},
			want: map[string]string{
				"foo": "",
			},
		},
		{
			name: "normal-2",
			args: args{
				values: []string{"foo=bar=baz"},
			},
			want: map[string]string{
				"foo": "bar=baz",
			},
		},
		{
			name: "empty",
			args: args{
				values: []string{},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := utils.KeyValueStringsToMap(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertKVStringsToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
