package templates

const fileTpl = `// Code generated by protoc-gen-zap
// source: {{ .InputPath }}
// DO NOT EDIT!!!
package {{ go_package_name .Package }}
import (
	"strconv"
	"go.uber.org/zap/zapcore"
	"github.com/bigbluedisco/protoc-gen-zap/utils"
)

// ensure the imports are used
var (
	_ = (*utils.Interfaces)(nil)
	_ = zapcore.InfoLevel
	_ = strconv.IntSize
)

{{ range .AllMessages }}
	{{ template "message" . }}
{{ end }}
`
