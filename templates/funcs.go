package templates

import (
	"bytes"
	"fmt"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

// Register adds a render function to the template
func Register(tpl *template.Template) {
	tpl.Funcs(map[string]interface{}{
		"render": render,
	})

	template.Must(tpl.Parse(fileTpl))
	template.Must(tpl.New("message").Parse(messageTpl))
}

func getter(n pgs.Name, t pgs.FieldType) string {
	s := fmt.Sprintf("c.Get%s()", n.UpperCamelCase())

	if t.IsEnum() || wellKnowType(t.Name().String()) {
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

func isSimple(t string) bool {
	switch t {
	case "float64",
		"float32",
		"int32",
		"int64",
		"uint32",
		"uint64",
		"bool",
		"string",
		"[]byte":
		return true
	}
	return false
}

func simpleAddFunc(n pgs.Name, t pgs.FieldType) string {

	if t.IsEnum() || wellKnowType(t.Name().String()) {
		return "AddString"
	}

	if t.IsEmbed() {
		return "AddObject"
	}

	switch t.Name() {
	case "float64":
		// proto: double
		return "AddFloat64"
	case "float32":
		// proto: float
		return "AddFloat32"
	case "int32":
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "AddInt32"
	case "int64":
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "AddInt64"
	case "uint32":
		// proto: uint32
		// proto: fixed32
		return "AddUint32"
	case "uint64":
		// proto: uint64
		// proto: fixed64
		return "AddUint64"
	case "bool":
		// proto: bool
		return "AddBool"
	case "string":
		// proto: string
		return "AddString"
	case "[]byte":
		// proto: bytes
		return "AddBinary"
	}

	return "AddReflected"
}

func arrayFunc(typ pgs.FieldType) string {

	t := typ.Element()

	switch t.Name() {
	case "float64":
		// proto: double
		return "Float64s"
	case "float32":
		// proto: float
		return "Float32s"
	case "int32":
		// proto: int32
		// proto: sint32
		// proto: sfixed32
		return "Int32s"
	case "int64":
		// proto: int64
		// proto: sint64
		// proto: sfixed64
		return "Int64s"
	case "uint32":
		// proto: uint32
		// proto: fixed32
		return "Uint32s"
	case "uint64":
		// proto: uint64
		// proto: fixed64
		return "Uint64s"
	case "bool":
		// proto: bool
		return "Bools"
	case "[]byte":
		// proto: bytes
		return "ByteStringsArray"
	}

	// proto: string
	return "StringArray"
}

func wellKnowType(t string) bool {

	switch t {
	case "*timestamp.Timestamp",
		"*empty.Empty",
		"*any.Any",
		"*duration.Duration",
		"*struct.Struct":
		return true
	}

	if len(t) < 11 {
		return false
	}

	// wrappers
	if t[:10] == "*wrappers." {
		return true
	}

	return false
}

const oneoftpl = `
if %s != nil {
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
	return ArrayData{
		SliceName: "s" + lowerString(4),
		IndexName: "i" + lowerString(4),
		ItemName:  "o" + lowerString(4),
		SliceType: sliceType,
		Getter:    getter,
		Key:       key,
	}
}

const arrayTpl = `
{{ .SliceName }} := make(utils.{{ .SliceType }}, len({{ .Getter }}))
for {{ .IndexName }}, {{ .ItemName }} := range {{ .Getter }} {
	{{ .SliceName }}[{{ .IndexName }}] = {{ .ItemName }}
}
o.AddArray("{{ .Key }}", {{ .SliceName }})
`

func render(f pgs.Field) string {
	t := f.Type()
	n := f.Name()

	var s string

	// repeated
	if t.IsRepeated() {

		if t.Element().IsEnum() || wellKnowType(t.Element().Name().String()) {

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

		} else if isSimple(t.Element().Name().String()) {

			s = fmt.Sprintf(`o.AddArray("%s", utils.%s(%s))`, name(f), arrayFunc(t), getter(n, t))

		} else {
			d := newArrayData("Interfaces", getter(n, t), name(f))
			tpl := template.New("interfaces")
			template.Must(tpl.Parse(arrayTpl))
			bb := bytes.NewBufferString("")
			_ = tpl.Execute(bb, d)

			s = bb.String()
		}
	} else {
		s = fmt.Sprintf(`o.%s("%s", %s)`, simpleAddFunc(n, t), name(f), getter(n, t))
	}

	// of oneof wrap in if != nil
	if f.OneOf() != nil {
		s = fmt.Sprintf(oneoftpl, getter(n, t), s)
	}

	return s
}
