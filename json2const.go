package json2enum

import (
	"bytes"
	"io"
	"strings"

	"github.com/antchfx/jsonquery"
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
func Convert(r io.Reader) (io.Reader, error) {
	return NewWithDefaultSetting().Convert(r)
}

func ConvertFromBytes(data []byte) (io.Reader, error) {
	return NewWithDefaultSetting().ConvertFromBytes(data)
}

func ConvertFromString(data string) (io.Reader, error) {
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

func (c Converter) Convert(r io.Reader) (io.Reader, error) {
	doc, err := jsonquery.Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error when load data")
	}
	list := jsonquery.Find(doc, c.PathToArray)
	var b bytes.Buffer
	for _, item := range list {
		if item == nil {
			continue
		}
		_, err := b.WriteString(item.InnerText())
		if err != nil {
			return nil, errors.Wrap(err, "Can't write response")
		}
	}
	return bytes.NewReader(b.Bytes()), nil
}

func (c Converter) ConvertFromBytes(data []byte) (io.Reader, error) {
	return c.Convert(bytes.NewBuffer(data))
}

func (c Converter) ConvertFromString(data string) (io.Reader, error) {
	return c.Convert(strings.NewReader(data))
}
