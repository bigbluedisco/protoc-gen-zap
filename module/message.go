package module

import (
	"bytes"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const wellKnowTypeImportPath = "google.golang.org/protobuf/types/known"

func generateMessage(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P("// MarshalLogObject makes ", message.GoIdent, " implement zap.ObjectMarshaler")
	g.P("func (c *", message.GoIdent, ") MarshalLogObject(o ", g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "ObjectEncoder",
		GoImportPath: "go.uber.org/zap/zapcore",
	}), ") error {")
	g.P("if c == nil {")
	g.P("return nil")
	g.P("}")
	g.P()

	for _, field := range message.Fields {
		if err := generateField(g, field); err != nil {
			return err
		}
	}

	g.P("return nil")
	g.P("}")
	g.P()

	return nil
}

func generateField(g *protogen.GeneratedFile, field *protogen.Field) error {
	// if oneof wrap in <if not empty>
	// not required if already wrapped (for embed types)
	var alreadyCheckedNil bool
	if field.Oneof != nil {
		g.P("if ", getter(field, false), " != ", zeroValue(field.Desc.Kind()), " {")
		alreadyCheckedNil = true
	}

	switch {
	case field.Desc.IsList():

		if field.Enum != nil || (field.Message != nil && strings.HasPrefix(string(field.Message.GoIdent.GoImportPath), wellKnowTypeImportPath)) {
			d := newArrayData("Stringers", getter(field, false), fieldName(field))
			if s, err := d.Render(g); err != nil {
				return err
			} else {
				g.P(s)
			}
		} else if field.Message != nil {
			d := newArrayData("Objects", getter(field, false), fieldName(field))
			if s, err := d.Render(g); err != nil {
				return err
			} else {
				g.P(s)
			}
		} else if isSimple(field.Desc.Kind()) {
			g.P("o.AddArray(\"", fieldName(field), "\", ", g.QualifiedGoIdent(protogen.GoIdent{
				GoImportPath: "github.com/bigbluedisco/protoc-gen-zap/utils",
				GoName:       arrayFunc(field.Desc.Kind()),
			}), "(", getter(field, false), "))")
		} else {
			d := newArrayData("Interfaces", getter(field, false), fieldName(field))
			if s, err := d.Render(g); err != nil {
				return err
			} else {
				g.P(s)
			}
		}

	case field.Message != nil && field.Desc.MapKey() == nil:

		if strings.HasPrefix(string(field.Message.GoIdent.GoImportPath), wellKnowTypeImportPath) {

			if !alreadyCheckedNil {
				g.P("if ", getter(field, true), " != nil {")
			}
			g.P("o.AddString(\"", fieldName(field), "\", ", getter(field, true), ".String())")
			if !alreadyCheckedNil {
				g.P("}")
			}
		} else {
			if !alreadyCheckedNil {
				g.P("if ", getter(field, true), " != nil {")
			}
			g.P("o.AddObject(\"", fieldName(field), "\", ", getter(field, true), ")")
			if !alreadyCheckedNil {
				g.P("}")
			}
		}

	default:
		g.P("o.", simpleAddFunc(field.Desc.Kind()), "(\"", fieldName(field), "\", ", getter(field, true), ")")
	}

	if field.Oneof != nil {
		g.P("}")
	}
	g.P("")

	return nil
}

func fieldName(field *protogen.Field) string {
	if field.Oneof != nil {
		return string(field.Oneof.Desc.Name())
	}
	return string(field.Desc.Name())
}

func getter(field *protogen.Field, enumToString bool) string {
	if field.Enum != nil && enumToString {
		return "c.Get" + field.GoName + "().String()"
	}
	return "c.Get" + field.GoName + "()"
}

func newArrayData(sliceType string, getter string, key string) arrayData {
	name := strings.Replace(key, "_", "", -1) + strings.ToLower(sliceType)
	return arrayData{
		SliceName: name,
		IndexName: name + "k",
		ItemName:  name + "v",
		SliceType: sliceType,
		Getter:    getter,
		Key:       key,
	}
}

type arrayData struct {
	SliceName string
	IndexName string
	ItemName  string
	SliceType string
	Getter    string
	Key       string
}

func (d arrayData) Render(g *protogen.GeneratedFile) (string, error) {
	d.SliceType = g.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/bigbluedisco/protoc-gen-zap/utils",
		GoName:       d.SliceType,
	})

	var b bytes.Buffer
	if err := arrayTpl.Execute(&b, d); err != nil {
		return "", err
	}
	return b.String(), nil
}

const arrayTplString = `
{{ .SliceName }}Length := len({{ .Getter }})
o.AddInt("{{ .Key }}_length", {{ .SliceName }}Length)

if {{ .SliceName }}Length > 100 {
	{{ .SliceName }}Length = 100
}

{{ .SliceName }} := make({{ .SliceType }}, {{ .SliceName }}Length)
for i := 0; i < {{ .SliceName }}Length; i++ {
	{{ .SliceName }}[i] = {{ .Getter }}[i]
}

if {{ .SliceName }}Length == 100 {
	o.AddArray("first_100_{{ .Key }}", {{ .SliceName }})
} else {
	o.AddArray("{{ .Key }}", {{ .SliceName }})
}
`

var arrayTpl = template.Must(template.New("array").Parse(arrayTplString))

func isSimple(t protoreflect.Kind) bool {
	switch t {
	case protoreflect.DoubleKind,
		protoreflect.FloatKind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.BoolKind,
		protoreflect.StringKind,
		protoreflect.BytesKind:
		return true
	}
	return false
}

func arrayFunc(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.DoubleKind:
		// proto: double
		return "Float64s"
	case protoreflect.FloatKind:
		// proto: float
		return "Float32s"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "Int32s"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "Int64s"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		// proto: uint32
		// proto: fixed32
		return "Uint32s"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		// proto: uint64
		// proto: fixed64
		return "Uint64s"
	case protoreflect.BoolKind:
		// proto: bool
		return "Bools"
	case protoreflect.BytesKind:
		// proto: bytes
		return "ByteStringsArray"
	}

	// proto: string
	return "StringArray"
}

func simpleAddFunc(t protoreflect.Kind) string {
	switch t {
	case protoreflect.EnumKind:
		return "AddString"
	case protoreflect.DoubleKind:
		// proto: double
		return "AddFloat64"
	case protoreflect.FloatKind:
		// proto: float
		return "AddFloat32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "AddInt32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "AddInt64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		// proto: uint32
		// proto: fixed32
		return "AddUint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		// proto: uint64
		// proto: fixed64
		return "AddUint64"
	case protoreflect.BoolKind:
		// proto: bool
		return "AddBool"
	case protoreflect.StringKind:
		// proto: string
		return "AddString"
	case protoreflect.BytesKind:
		// proto: bytes
		return "AddBinary"
	}

	return "AddReflected"
}

func zeroValue(typ protoreflect.Kind) string {
	switch typ {
	case protoreflect.DoubleKind,
		protoreflect.FloatKind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.EnumKind:
		return "0"

	case protoreflect.BoolKind:
		return "false"

	case protoreflect.StringKind:
		return `""`

	case protoreflect.BytesKind, protoreflect.MessageKind:
		return "nil"

	default:
		return "nil"
	}
}
