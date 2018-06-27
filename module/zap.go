package module

import (
	"github.com/bigbluedisco/protoc-gen-zap/templates"
	pgs "github.com/lyft/protoc-gen-star"
)

const (
	zapName   = "zap"
	langParam = "lang"
)

type Module struct {
	*pgs.ModuleBase
}

func Zap() Module { return Module{&pgs.ModuleBase{}} }

func (m Module) Name() string { return zapName }

func (m Module) Execute(target pgs.Package, packages map[string]pgs.Package) []pgs.Artifact {
	lang := m.Parameters().Str(langParam)
	m.Assert(lang != "", "`lang` parameter must be set")
	tpl := templates.Template(lang)
	m.Assert(tpl != nil, "could not find template for `lang`: ", lang)

	for _, f := range target.Files() {
		m.Push(f.Name().String())

		m.AddGeneratorTemplateFile(
			f.OutputPath().SetExt(".zap."+tpl.Name()).String(),
			tpl,
			f,
		)

		m.Pop()
	}

	return m.Artifacts()
}

var _ pgs.Module = Module{}
