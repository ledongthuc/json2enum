package json2enum

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestConverter_ConvertFromString(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name      string
		converter Converter
		args      args
		want      string
		hasError  bool
	}{
		{
			name: "root path",
			converter: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "Cate",
				TypeName:    "MyCategory",
				PackageName: "mypackage",
			},
			args: args{
				data: `["stock", "bond", "real estate", "bitcoin & digital coins"]`,
			},
			want:     `abc`,
			hasError: false,
		},
		{
			name: "object with field",
			converter: Converter{
				PathToArray: "#.name",
				EnumPrefix:  "Cate",
				TypeName:    "MyCategory",
				PackageName: "mypackage",
			},
			args: args{
				data: `[{"name":"stock"},{"name":"bond"},{"name":"real estate"},{"name":"bitcoin & digital coins"}]`,
			},
			want:     `abc`,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := tt.converter.ConvertFromString(tt.args.data)
			fmt.Println(err != nil, tt.hasError)
			if err != nil && !tt.hasError {
				t.Errorf("Converter.ConvertFromString() has error = %v, want %v. Error content: (%v)", err != nil, tt.hasError, err)
				return
			}
			got, errf := ioutil.ReadAll(reader)
			if errf != nil {
				t.Errorf("Can't read data at: %s", errf)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Converter.ConvertFromString() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestConverter_SetDefault(t *testing.T) {
	type fields struct {
		PathToArray string
		EnumPrefix  string
		TypeName    string
		PackageName string
	}
	tests := []struct {
		name      string
		converter Converter
		want      Converter
	}{
		{
			name:      "Empty",
			converter: Converter{},
			want: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "",
				TypeName:    "MyType",
				PackageName: "json2enum",
			},
		},
		{
			name: "Custom",
			converter: Converter{
				PathToArray: "department.students",
				EnumPrefix:  "DepartmentType",
				TypeName:    "DepartmentType",
				PackageName: "json2enum",
			},
			want: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "",
				TypeName:    "MyType",
				PackageName: "json2enum",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.converter.SetDefault()
			if !reflect.DeepEqual(tt.converter, tt.want) {
				t.Errorf("Converter.SetDefault() got = %v, want %v", tt.converter, tt.want)
			}
		})
	}
}

func TestConverter_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		converter Converter
		valid     bool
		wantErr   error
	}{
		{
			name: "Successful",
			converter: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "MyType",
				TypeName:    "MyType",
				PackageName: "json2enum",
			},
			valid:   true,
			wantErr: nil,
		},
		{
			name: "Empty PathToArray",
			converter: Converter{
				PathToArray: "",
				EnumPrefix:  "MyType",
				TypeName:    "MyType",
				PackageName: "json2enum",
			},
			valid:   false,
			wantErr: ErrorInvalidArrayToPath,
		},
		{
			name: "Empty EnumPrefix",
			converter: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "",
				TypeName:    "MyType",
				PackageName: "json2enum",
			},
			valid:   true,
			wantErr: nil,
		},
		{
			name: "Empty TypeName",
			converter: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "MyType",
				TypeName:    "",
				PackageName: "json2enum",
			},
			valid:   false,
			wantErr: ErrorInvalidTypeName,
		},
		{
			name: "Empty PackageName",
			converter: Converter{
				PathToArray: "#()#",
				EnumPrefix:  "MyType",
				TypeName:    "MyType",
				PackageName: "",
			},
			valid:   false,
			wantErr: ErrorInvalidPackageName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := tt.converter.IsValid()
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Converter.IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.valid {
				t.Errorf("Converter.IsValid() = %v, want %v", valid, tt.valid)
			}
		})
	}
}
