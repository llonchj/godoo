package generator

import (
	"go/build"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
)

var apiTmpl *template.Template
var typesTmpl *template.Template

func init() {
	var err error

	funcMap := template.FuncMap{
		"funcName": strcase.ToCamel,
	}

	apiTmpl, err = template.New("api.tmpl").Funcs(funcMap).Parse(FSMustString(false, "/tmpl/api.tmpl"))
	if err != nil {
		panic(err)
	}

	typesTmpl, err = template.New("types.tmpl").Parse(FSMustString(false, "/tmpl/types.tmpl"))
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
