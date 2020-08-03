package templates

import (
	"text/template"

	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func Template(ctx pgsgo.Context, lang string) *template.Template {

	// Only go is supported
	if lang != "go" {
		return nil
	}

	tpl := template.New(lang)
	Register(ctx, tpl)
	return tpl
}
