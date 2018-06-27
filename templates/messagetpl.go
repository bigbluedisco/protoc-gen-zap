package templates

const messageTpl = `

// MarshalLogObject makes {{ .TypeName }} implement zap.ObjectMarshaler
func (c {{ .TypeName }}) MarshalLogObject(o zapcore.ObjectEncoder) error {

	{{ range .Fields }}
		
		{{ render . }}
		
	{{ end }}

	return nil
}
`
