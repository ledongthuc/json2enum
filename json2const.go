package json2enum

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// [ ] take Json array
// [ ] array path can be define path
// [ ] loop through array
// [ ] get field and break it as enum
// [ ] Support function to get back json object wrapper
// [ ] Support string and object
// [ ] Support enum Prefix
// [ ] Support package custom

// Lib: https://github.com/antchfx/jsonquery
// Lib: https://github.com/stretchr/objx

var (
	ErrorCantReadData = errors.New("Can't read stream json data")
)

// Convert json to enum bases on reader stream. Returned reader's always nil when has error.
func Convert(r io.Reader) (io.Reader, err) {
	return NewWithDefaultSetting().Convert(r)
}

func ConvertFromBytes(data []byte) (io.Reader, err) {
	return NewWithDefaultSetting().ConvertFromBytes(data)
}

func ConvertFromString(data string) (io.Reader, err) {
	return NewWithDefaultSetting().ConvertFromString(data)
}

type Converter struct {
	PathToArray       string
	EnumPrefix        string
	CustomPackageName string
}

func New() Converter {
	var c Converter
	return c
}

func NewWithDefaultSetting() Converter {
	c := New()
	c.SetDefault()
	return c
}

func (c Converter) SetDefault() {
	c.PathToArray = "//"
	c.EnumPrefix = ""
	c.CustomPackageName = "json2enum"
}

func (c Converter) Convert(r io.Reader) (io.Reader, err) {
	data, err := ioutil.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, ErrorCantReadData)
	}
	return c.ConvertFromBytes(data)
}

func (c Converter) ConvertFromBytes(data []byte) (io.Reader, err) {
	return c.ConvertFromString(data)
}

func (c Converter) ConvertFromString(data string) (io.Reader, err) {
	return nil, errors.New("Under construction!")
}
