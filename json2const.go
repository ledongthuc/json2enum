package json2enum

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// [ ] take Json array
// [ ] array path can be define path
// [ ] loop through array
// [ ] get field and break it as enum
// [ ] Support function to get back json object wrapper
// [ ] Support string and object
// [ ] Support enum Prefix
// [ ] Support package custom
// [ ] Format output
// [ ] Support special unicode characters like !@#$ to texable Variable name

// Lib: https://github.com/antchfx/jsonquery
// Lib: https://github.com/stretchr/objx
// Lib: https://golang.org/src/go/format/format.go?s=3141:3180#L81

var (
	ErrorCantReadData       = errors.New("Can't read stream json data")
	ErrorInvalidArrayToPath = errors.New("Array's path is not valid value. If json data is array of string, you can use '#()#'")
	ErrorInvalidTypeName    = errors.New("Type name is not valid.")
	ErrorInvalidPackageName = errors.New("Package name is not valid.")
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
	PathToArray string
	EnumPrefix  string
	TypeName    string
	PackageName string
}

// New will create a instance of Converter, that's useful for adding custom value to generate.
func New() Converter {
	var c Converter
	return c
}

// NewWithDefaultSetting will create a instance of Converter, that's useful for adding custom value to generate.
func NewWithDefaultSetting() Converter {
	c := New()
	c.SetDefault()
	return c
}

// SetDefault set input data with default values. It's used for run without entering any custom configuration to generate source code.
func (c *Converter) SetDefault() {
	c.PathToArray = "#()#"
	c.EnumPrefix = ""
	c.TypeName = "MyType"
	c.PackageName = "json2enum"
}

// IsValid validates input data before converting it to source file.
func (c *Converter) IsValid() (bool, error) {
	if c.PathToArray == "" {
		return false, ErrorInvalidArrayToPath
	}
	if c.TypeName == "" {
		return false, ErrorInvalidTypeName
	}
	if c.PackageName == "" {
		return false, ErrorInvalidPackageName
	}
	return true, nil
}

func (c Converter) Convert(r io.Reader) (io.Reader, error) {
	if valid, err := c.IsValid(); !valid {
		return nil, err
	}
	json, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error when load data")
	}

	// Prepare parameters
	tmplParameters := TemplateParameters{
		PackageName: c.PackageName,
		Type:        c.TypeName,
	}
	tmplParameters.GenerateTypeSingular()
	tmplParameters.GenerateTypePlural()
	fields := gjson.GetBytes(json, c.PathToArray)
	for _, item := range fields.Array() {
		if item.Type == gjson.Null {
			continue
		}

		tmplParameters.Fields = append(tmplParameters.Fields, TemplateField{
			Name:  fmt.Sprintf("%s%s", c.EnumPrefix, strcase.ToCamel(item.String())), // TODO: Support field name valid
			Value: item.String(),
		})
	}

	// Generate source code
	tmpl, err := template.New("template").Parse(json2constTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "Error when parsing template")
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, tmplParameters)
	if err != nil {
		return nil, errors.Wrap(err, "Error when generate source code")
	}
	a, err := format.Source(result.Bytes())
	if err != nil {
		fmt.Println("DEBUG: ", string(result.Bytes()))
		return nil, errors.Wrap(err, "Error when format generated code")
	}
	fmt.Println("DEBUG: ", string(a))
	return nil, errors.New("Test")
	return result, nil
}

func (c Converter) ConvertFromBytes(data []byte) (io.Reader, error) {
	return c.Convert(bytes.NewBuffer(data))
}

func (c Converter) ConvertFromString(data string) (io.Reader, error) {
	return c.Convert(strings.NewReader(data))
}
