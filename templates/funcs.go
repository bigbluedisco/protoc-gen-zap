package templates

import (
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

	if t.IsEmbed() {
		return "Objects"
	}

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
	case "string":
		// proto: string
		return "StringArray"
	case "[]byte":
		// proto: bytes
		return "ByteStringsArray"
	}

	return "Interfaces"
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

const stringerArrayTpl = `
gd2 := make(utils.Stringers, len(%s))
for i12, afd3 := range %s {
	gd2[i12] = afd3
}
o.AddArray("%s", utils.Stringers(gd2))
`

func render(f pgs.Field) string {
	t := f.Type()
	n := f.Name()

	var s string

	// repeated
	if t.IsRepeated() {
		if t.Element().IsEnum() || wellKnowType(t.Element().Name().String()) {
			s = fmt.Sprintf(stringerArrayTpl, getter(n, t), getter(n, t), name(f))
		} else {
			s = fmt.Sprintf(`o.AddArray("%s", utils.%s(%s))`, name(f), arrayFunc(t), getter(n, t))
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
