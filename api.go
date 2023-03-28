package generator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func GenerateBaseAPI(pkg, path, basePath string) error {
	for _, file := range []string{"client.go"} {
		b, err := content.ReadFile("api/" + file)
		if err != nil {
			return err
		}
		s := strings.Replace(string(b), "package api", "package "+pkg, 1)
		s = strings.Replace(s,
			"\"github.com/llonchj/godoo/types\"",
			"\""+filepath.Join(path, pkg, "types")+"\"", 1)
		if err := ioutil.WriteFile(filepath.Join(basePath, file), []byte(s), 0644); err != nil {
			return err
		}
	}
	return nil
}

func GenerateAPI(pkg, path, basePath, model string) error {
	snakeModel := strings.Replace(model, ".", "_", -1)
	var outTpl bytes.Buffer
	args := map[string]string{
		"Package": pkg,
		"Path":    path,
		"Model":   snakeModel,
	}
	err := apiTmpl.Execute(&outTpl, args)
	if err != nil {
		return err
	}

	b, err := format.Source(outTpl.Bytes())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(basePath, snakeModel+"_gen.go"), b, 0644)
}
