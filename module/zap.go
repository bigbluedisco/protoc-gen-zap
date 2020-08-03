package module

import (
	"github.com/bigbluedisco/protoc-gen-zap/templates"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const (
	zapName   = "zap"
	langParam = "lang"
)

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func Zap() pgs.Module { return &Module{ModuleBase: &pgs.ModuleBase{}} }

func (m Module) Name() string { return zapName }

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m Module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	lang := m.Parameters().Str(langParam)
	m.Assert(lang != "", "`lang` parameter must be set")
	tpl := templates.Template(m.ctx, lang)
	m.Assert(tpl != nil, "could not find template for `lang`: ", lang)

	for _, f := range targets {
		m.Push(f.Name().String())

		m.AddGeneratorTemplateFile(
			m.ctx.OutputPath(f).SetExt(".zap."+tpl.Name()).String(),
			tpl,
			f,
		)

		m.Pop()
	}

	return m.Artifacts()
}

var _ pgs.Module = (*Module)(nil)
