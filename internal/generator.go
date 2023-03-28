package generator

import (
	"embed"
	"go/build"
	"os"
	"text/template"

	"github.com/serenize/snaker"
)

//go:embed tmpl/* types/* api/client.go
var content embed.FS

var apiTmpl *template.Template
var typesTmpl *template.Template

func init() {
	var err error

	funcMap := template.FuncMap{
		"funcName": snaker.SnakeToCamel,
		// "funcName": strcase.ToCamel,
	}

	api, err := content.ReadFile("tmpl/api.tmpl")
	if err != nil {
		panic(err)
	}
	types, err := content.ReadFile("tmpl/types.tmpl")
	if err != nil {
		panic(err)
	}

	apiTmpl, err = template.New("api.tmpl").Funcs(funcMap).Parse(string(api))
	if err != nil {
		panic(err)
	}

	typesTmpl, err = template.New("types.tmpl").Parse(string(types))
	if err != nil {
		panic(err)
	}
}

func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}
