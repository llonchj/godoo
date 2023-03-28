package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/serenize/snaker"
)

var ErrMissingFieldType = errors.New("missing type")

var convertTypes = map[string]string{
	"datetime":           "time.Time",
	"date":               "time.Time",
	"monetary":           "float64",
	"char":               "string",
	"many2one":           "Many2One",
	"many2one_reference": "int64",
	"many2many":          "[]int64",
	"one2many":           "[]int64",
	"integer":            "int64",
	"boolean":            "bool",
	"text":               "string",
	"selection":          "interface{}",
	"float":              "float64",
	"binary":             "string",
	"html":               "string",
	"reference":          "string",
}

var convertNilTypes = map[string]string{
	"datetime":           "interface{}",
	"date":               "interface{}",
	"monetary":           "interface{}",
	"char":               "interface{}",
	"many2one":           "interface{}",
	"many2one_reference": "interface{}",
	"many2many":          "interface{}",
	"one2many":           "interface{}",
	"integer":            "interface{}",
	"boolean":            "bool",
	"text":               "interface{}",
	"selection":          "interface{}",
	"float":              "interface{}",
	"binary":             "interface{}",
	"html":               "interface{}",
	"reference":          "interface{}",
}

type ModelType struct {
	ModelName      string
	CamelModelName string
	Fields         []Field
	Time           bool
}

type Field struct {
	Name      string
	SnakeName string
	Type      string
	NilType   string
}

func GenerateBaseTypes(pkg, path, basePath string) error {
	b, err := content.ReadFile("types/types.go")
	if err != nil {
		return err
	}

	ioutil.WriteFile(filepath.Join(basePath, "types.go"), b, 0644)
	return nil
}

func GenerateTypes(pkg, path, basePath, model string, fields map[string]string) error {
	snakeModel := strings.Replace(model, ".", "_", -1)
	modelType := ModelType{ModelName: model, CamelModelName: snaker.SnakeToCamel(snakeModel)}
	for fieldName, fieldType := range fields {
		convertType, found := convertTypes[fieldType]
		if !found {
			return fmt.Errorf("%w: %s", ErrMissingFieldType, fieldType)
		}
		if convertType == "time.Time" {
			modelType.Time = true
		}
		t, found := convertNilTypes[fieldType]
		if !found {
			return fmt.Errorf("%w nil: %s", ErrMissingFieldType, fieldType)
		}

		f := Field{Name: snaker.SnakeToCamel(fieldName), SnakeName: fieldName, Type: convertType, NilType: t}
		modelType.Fields = append(modelType.Fields, f)
	}

	var outTpl bytes.Buffer
	err := typesTmpl.Execute(&outTpl, modelType)
	if err != nil {
		return err
	}

	b, err := format.Source(outTpl.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(basePath, snakeModel+"_gen.go"), b, 0644)
}
