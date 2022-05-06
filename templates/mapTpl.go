package templates

import (
	privacy "github.com/bigbluedisco/protoc-gen-privacy"
	pgs "github.com/lyft/protoc-gen-star"
)

type mapData struct {
	MapName   string
	Getter    string
	AddFunc   string
	FormatKey string
	IsStars   bool
}

func newMapData(f pgs.Field, obsType privacy.Rule) mapData {
	formatKey := "key"
	keyType := f.Type().Key().ProtoType()
	if isSimple(keyType) && keyType != pgs.StringT {
		formatKey = "strconv.Itoa(int(key))"
	}

	return mapData{
		MapName:   name(f),
		Getter:    getter(f.Name(), f.Type()),
		FormatKey: formatKey,
		IsStars:   obsType == privacy.Rule_STARS,
		AddFunc:   simpleAddFunc(f.Type().Element()),
	}
}

const mapTpl = `
{{ if .IsStars }}
for key := range {{ .Getter }} {
		o.AddString("{{ .MapName }}_" + key, "***")
}
{{ else }}
for key, value := range {{ .Getter }} {
		o.{{ .AddFunc }}("{{ .MapName }}_" + {{ .FormatKey }}, value)
}
{{ end }}
`
