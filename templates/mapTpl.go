package templates

import (
	"github.com/bigbluedisco/protoc-gen-zap/zap"
	pgs "github.com/lyft/protoc-gen-star"
)

type mapData struct {
	MapName   string
	Getter    string
	AddFunc   string
	FormatKey string
	IsStars   bool
}

func newMapData(f pgs.Field, obsType zap.ObfuscationType) mapData {
	formatKey := "key"
	keyType := f.Type().Key().ProtoType()
	if isSimple(keyType) && keyType != pgs.StringT {
		formatKey = "strconv.Itoa(int(key))"
	}

	return mapData{
		MapName:   name(f),
		Getter:    getter(f.Name(), f.Type()),
		FormatKey: formatKey,
		IsStars:   obsType == zap.ObfuscationType_STARS,
		AddFunc:   simpleAddFunc(f.Type().Element()),
	}
}

const mapTpl = `
{{ if .IsStars }}
for key := range {{ .Getter }} {
		o.AddString("{{ .MapName }}_" + key, "***")
		continue
}
{{ else }}
for key, value := range {{ .Getter }} {
		o.{{ .AddFunc }}("{{ .MapName }}_" + {{ .FormatKey }}, value)
}
{{ end }}
`
