package templates

const messageTpl = `

// MarshalLogObject makes {{ go_type_name . }} implement zap.ObjectMarshaler
func (c {{ go_type_name . }}) MarshalLogObject(o zapcore.ObjectEncoder) error {

	{{ range .Fields }}
		
		{{ render . }}
		
	{{ end }}

	return nil
}
`
