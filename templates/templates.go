package templates

import "text/template"

func Template(lang string) *template.Template {

	// Only go is supported
	if lang != "go" {
		return nil
	}

	tpl := template.New(lang)
	Register(tpl)
	return tpl
}
