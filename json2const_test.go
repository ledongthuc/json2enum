package json2enum

import (
	"reflect"
	"testing"
)

func TestConverter_ConvertFromString(t *testing.T) {
	type fields struct {
		PathToArray       string
		EnumPrefix        string
		CustomPackageName string
	}
	type args struct {
		data string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		hasError bool
	}{
		{
			name: "root path",
			fields: fields{
				PathToArray:       "//",
				EnumPrefix:        "",
				CustomPackageName: "test",
			},
			args: args{
				data: `["a", "b", "c"]`,
			},
			want:     ``,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Converter{
				PathToArray:       tt.fields.PathToArray,
				EnumPrefix:        tt.fields.EnumPrefix,
				CustomPackageName: tt.fields.CustomPackageName,
			}
			got, err := c.ConvertFromString(tt.args.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Converter.ConvertFromString() got = %v, want %v", got, tt.want)
			}
			if err == nil && !tt.hasError {
				t.Errorf("Converter.ConvertFromString() has error = %v, want %v", err != nil, tt.hasError)
			}
		})
	}
}
