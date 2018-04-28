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

	apiTmpl, err = template.New("api.tmpl.go").Funcs(funcMap).Parse(FSMustString(false, "/tmpl/api.tmpl.go"))
	if err != nil {
		panic(err)
	}

	typesTmpl, err = template.New("types.tmpl.go").Parse(FSMustString(false, "/tmpl/types.tmpl.go"))
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
