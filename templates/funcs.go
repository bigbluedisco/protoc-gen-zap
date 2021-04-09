package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// Register adds a render function to the template
func Register(ctx pgsgo.Context, tpl *template.Template) {
	tpl.Funcs(map[string]interface{}{
		"render": render,
		"go_package_name": func(p pgs.Package) string {
			return ctx.PackageName(p).String()
		},
		"go_type_name": func(m pgs.Message) string {
			return ctx.Name(m).String()
		},
	})

	template.Must(tpl.Parse(fileTpl))
	template.Must(tpl.New("message").Parse(messageTpl))
}

func getter(n pgs.Name, t pgs.FieldType) string {
	s := fmt.Sprintf("c.Get%s()", n.UpperCamelCase())

	if t.IsEnum() {
		return s + ".String()"
	}

	return s
}

func name(f pgs.Field) string {

	if f.InOneOf() {
		return f.OneOf().Name().String()
	}

	return f.Name().String()
}

func isSimple(t pgs.ProtoType) bool {
	switch t {
	case pgs.DoubleT,
		pgs.FloatT,
		pgs.Int32T, pgs.SInt32, pgs.SFixed32,
		pgs.Int64T, pgs.SInt64, pgs.SFixed64,
		pgs.UInt32T, pgs.Fixed32T,
		pgs.UInt64T, pgs.Fixed64T,
		pgs.BoolT,
		pgs.StringT,
		pgs.BytesT:
		return true
	}
	return false
}

func simpleAddFunc(n pgs.Name, t pgs.FieldType) string {

	switch t.ProtoType() {
	case pgs.EnumT:
		return "AddString"
	case pgs.DoubleT:
		// proto: double
		return "AddFloat64"
	case pgs.FloatT:
		// proto: float
		return "AddFloat32"
	case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "AddInt32"
	case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "AddInt64"
	case pgs.UInt32T, pgs.Fixed32T:
		// proto: uint32
		// proto: fixed32
		return "AddUint32"
	case pgs.UInt64T, pgs.Fixed64T:
		// proto: uint64
		// proto: fixed64
		return "AddUint64"
	case pgs.BoolT:
		// proto: bool
		return "AddBool"
	case pgs.StringT:
		// proto: string
		return "AddString"
	case pgs.BytesT:
		// proto: bytes
		return "AddBinary"
	}

	return "AddReflected"
}

func arrayFunc(typ pgs.FieldType) string {

	switch typ.Element().ProtoType() {
	case pgs.DoubleT:
		// proto: double
		return "Float64s"
	case pgs.FloatT:
		// proto: float
		return "Float32s"
	case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "Int32s"
	case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "Int64s"
	case pgs.UInt32T, pgs.Fixed32T:
		// proto: uint32
		// proto: fixed32
		return "Uint32s"
	case pgs.UInt64T, pgs.Fixed64T:
		// proto: uint64
		// proto: fixed64
		return "Uint64s"
	case pgs.BoolT:
		// proto: bool
		return "Bools"
	case pgs.BytesT:
		// proto: bytes
		return "ByteStringsArray"
	}

	// proto: string
	return "StringArray"
}

func zeroValue(typ pgs.FieldType) string {
	switch typ.ProtoType() {
	case pgs.DoubleT,
		pgs.FloatT,
		pgs.Int32T, pgs.SInt32, pgs.SFixed32,
		pgs.Int64T, pgs.SInt64, pgs.SFixed64,
		pgs.UInt32T, pgs.Fixed32T,
		pgs.UInt64T, pgs.Fixed64T, pgs.EnumT:
		return "0"

	case pgs.BoolT:
		return "false"

	case pgs.StringT:
		return `""`

	case pgs.BytesT, pgs.MessageT:
		return "nil"

	default:
		return "nil"
	}
}

const oneoftpl = `
if %s != %s {
	%s
}
`

// ArrayData used to fill arrayTpl
type ArrayData struct {
	SliceName string
	IndexName string
	ItemName  string
	SliceType string
	Getter    string
	Key       string
}

func newArrayData(sliceType string, getter string, key string) ArrayData {
	name := strings.Replace(key, "_", "", -1) + strings.ToLower(sliceType)
	return ArrayData{
		SliceName: name,
		IndexName: name + "k",
		ItemName:  name + "v",
		SliceType: sliceType,
		Getter:    getter,
		Key:       key,
	}
}

const arrayTpl = `
{{ .SliceName }}Length := len({{ .Getter }})
o.AddInt("{{ .Key }}_length", {{ .SliceName }}Length)

if {{ .SliceName }}Length > 100 {
	{{ .SliceName }}Length = 100
}

{{ .SliceName }} := make(utils.{{ .SliceType }}, {{ .SliceName }}Length)
for i := 0; i < {{ .SliceName }}Length; i++ {
	{{ .SliceName }}[i] = {{ .Getter }}[i]
}

if {{ .SliceName }}Length == 100 {
	o.AddArray("first_100_{{ .Key }}", {{ .SliceName }})
} else {
	o.AddArray("{{ .Key }}", {{ .SliceName }})
}
`

func render(f pgs.Field) string {
	t := f.Type()
	n := f.Name()

	var s string

	// repeated
	if t.IsRepeated() {

		if t.Element().IsEnum() || (t.Element().IsEmbed() && t.Element().Embed().IsWellKnown()) {
			d := newArrayData("Stringers", getter(n, t), name(f))
			tpl := template.New("stringers")
			template.Must(tpl.Parse(arrayTpl))
			bb := bytes.NewBufferString("")
			_ = tpl.Execute(bb, d)

			s = bb.String()

		} else if t.Element().IsEmbed() {
			d := newArrayData("Objects", getter(n, t), name(f))
			tpl := template.New("objects")
			template.Must(tpl.Parse(arrayTpl))
			bb := bytes.NewBufferString("")
			_ = tpl.Execute(bb, d)

			s = bb.String()

		} else if isSimple(t.Element().ProtoType()) {
			s = fmt.Sprintf(`o.AddArray("%s", utils.%s(%s))`, name(f), arrayFunc(t), getter(n, t))

		} else {
			d := newArrayData("Interfaces", getter(n, t), name(f))
			tpl := template.New("interfaces")
			template.Must(tpl.Parse(arrayTpl))
			bb := bytes.NewBufferString("")
			_ = tpl.Execute(bb, d)

			s = bb.String()
		}
	} else if t.IsEmbed() {
		if t.Embed().IsWellKnown() {
			s = fmt.Sprintf(`if %s != nil {
				o.AddString("%s", %s.String())
			}`, getter(n, t), name(f), getter(n, t))
		} else {
			s = fmt.Sprintf(`if %s != nil {
				o.AddObject("%s", %s)
			}`, getter(n, t), name(f), getter(n, t))
		}
	} else {
		s = fmt.Sprintf(`o.%s("%s", %s)`, simpleAddFunc(n, t), name(f), getter(n, t))
	}

	// if oneof wrap in <if not empty>
	// not required if already wrapped (for embed types)
	if f.OneOf() != nil && !strings.Contains(s, "!= nil") {
		s = fmt.Sprintf(oneoftpl, getter(n, t), zeroValue(t), s)
	}

	return s
}
