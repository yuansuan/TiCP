package template

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

func Render(text string, data interface{}) (string, error) {
	tpl, err := template.New("default").Parse(text)
	if err != nil {
		return "", errors.Wrap(err, "parse")
	}

	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, data); err != nil {
		return "", errors.Wrap(err, "render")
	}
	return buf.String(), nil
}
