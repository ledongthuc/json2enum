package json2enum

import "github.com/jinzhu/inflection"

const json2constTemplate = `package {{.PackageName}}

type {{.TypeSingular}} string
type {{.TypePlural}} []{{.TypeSingular}}

const (
	{{- range $index, $field := .Fields}}
		{{$field.Name}} = {{$field.Value | printf "\"%s\"" -}}
	{{end}}
)

var {{.TypePlural}}List = {{.TypePlural}}{
	{{- range $index, $field := .Fields}}
		{{$field.Name | printf "%s," -}}
	{{end}}
}

func (t {{.TypeSingular}}) ToPointer() *{{.TypeSingular}} {
	return &t
}

func (t *{{.TypeSingular}}) Value() {{.TypeSingular}} {
	if t == nil {
		return {{.TypeSingular}}{}
	}

	return *t
}
`

type TemplateParameters struct {
	PackageName  string
	Type         string
	TypeSingular string
	TypePlural   string
	Fields       TemplateFields
}

type TemplateField struct {
	Name  string
	Value string
}

type TemplateFields []TemplateField

func (t *TemplateParameters) GenerateTypeSingular() {
	if len(t.TypeSingular) > 0 {
		return
	}
	t.TypeSingular = inflection.Singular(t.Type)
}

func (t *TemplateParameters) GenerateTypePlural() {
	if len(t.TypePlural) > 0 {
		return
	}
	t.TypePlural = inflection.Plural(t.Type)
}
