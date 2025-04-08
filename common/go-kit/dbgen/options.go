package dbgen

import (
	"fmt"
	"strings"
	"text/template"
)

type Tags []*Tag

func (ts Tags) String() string {
	if len(ts) == 0 {
		return ""
	}
	var tagStrings []string
	for _, t := range ts {
		tagStrings = append(tagStrings, t.String())
	}

	return fmt.Sprintf("`%s`", strings.Join(tagStrings, " "))
}

func NewTags(tags ...*Tag) Tags {
	return tags
}

type Tag struct {
	key   string
	value string
}

func NewTag(key, value string) *Tag {
	return &Tag{
		key:   key,
		value: value,
	}
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:\"%s\"", t.key, t.value)
}

type ExtraField struct {
	Line          string
	ImportPackage string
}

type Options struct {
	GoPackage      string // main by default
	GeneratorName  string // dbgen by default
	DisableNull    bool
	Tables         []string                 // all tables if it's nil or empty
	OutputPathFunc func(table Table) string // os.Stdout if it's nil or returns empty string

	FieldTypeFunc           func(column ColumnSchema) (fieldType string, importPackage string, err error) // general field type if returns empty field type
	FieldNameFunc           func(column ColumnSchema) string
	TagsFunc                func(column ColumnSchema) Tags
	IgnoreColumnFunc        func(column ColumnSchema) bool
	StructureNameFunc       func(table Table) string
	TemplatePreparationFunc func(table Table) (interface{}, []string)
	Template                *template.Template
	EnableGoImports         bool
	ExtraFieldFunc          func(table Table) []*ExtraField
}

func (o *Options) merge(opt Options) {
	o.DisableNull = opt.DisableNull
	o.Tables = opt.Tables
	o.EnableGoImports = opt.EnableGoImports

	if opt.GoPackage != "" {
		o.GoPackage = opt.GoPackage
	}

	if opt.GeneratorName != "" {
		o.GeneratorName = opt.GeneratorName
	}

	if opt.OutputPathFunc != nil {
		o.OutputPathFunc = opt.OutputPathFunc
	}

	if opt.FieldTypeFunc != nil {
		o.FieldTypeFunc = opt.FieldTypeFunc
	}

	if opt.FieldNameFunc != nil {
		o.FieldNameFunc = opt.FieldNameFunc
	}

	if opt.StructureNameFunc != nil {
		o.StructureNameFunc = opt.StructureNameFunc
	}

	if opt.TagsFunc != nil {
		o.TagsFunc = opt.TagsFunc
	}

	if opt.TemplatePreparationFunc != nil {
		o.TemplatePreparationFunc = opt.TemplatePreparationFunc
	}

	if opt.Template != nil {
		o.Template = opt.Template
	}

	if opt.IgnoreColumnFunc != nil {
		o.IgnoreColumnFunc = opt.IgnoreColumnFunc
	}

	if opt.ExtraFieldFunc != nil {
		o.ExtraFieldFunc = opt.ExtraFieldFunc
	}
}

var defaultOptions = Options{
	GoPackage:               "main",
	GeneratorName:           "dbgen",
	FieldNameFunc:           func(column ColumnSchema) string { return snakeToCamel(column.ColumnName) },
	FieldTypeFunc:           func(column ColumnSchema) (string, string, error) { return column.goFieldType() },
	StructureNameFunc:       func(table Table) string { return snakeToCamel(table.Name) },
	IgnoreColumnFunc:        func(_ ColumnSchema) bool { return false },
	TemplatePreparationFunc: func(table Table) (interface{}, []string) { return &table, nil },
	Template: template.Must(template.New("tpl").Parse(`
// {{ .GoStructureName }}
type {{ .GoStructureName }} struct {
	{{range .Columns -}}
		{{ .GoFieldName }}  {{ .GoFieldType }}  {{ .GoTags }}
	{{end -}}

	{{range .ExtraFields}}
		{{ .Line -}}
	{{end}}
}
`)),
}

func DefaultOptions() Options {
	return defaultOptions
}

func snakeToCamel(s string) string {
	s = strings.Trim(s, " ")
	n := ""
	capNext := true
	for _, v := range s {
		if (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') {
			n += string(v)
		}

		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}

		capNext = v == '_' || v == '-'
	}

	return n
}
