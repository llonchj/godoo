package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	generator "github.com/llonchj/godoo/internal"

	"github.com/spf13/cobra"
)

func upsert(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("At least one argument needed")
		return
	}
	c, err := getSession(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := c.CompleteSession(); err != nil {
		fmt.Println(err.Error())
		return
	}

	if args[0] == "all" {
		args, err = c.GetAllModels()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	pkg, err := cmd.PersistentFlags().GetString("package")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	path, err := cmd.PersistentFlags().GetString("path")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dir := filepath.Join(path, pkg)
	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if strings.HasSuffix(wd, path) {
			dir = filepath.Join(wd, pkg)
		} else {
			dir = filepath.Join(wd, path, pkg)
		}
		if path == "" {
			_, path = filepath.Split(wd)
		}

	}
	// fmt.Printf("PKG: %s PATH: %s DIR: %s\n", pkg, path, dir)

	typesDir := filepath.Join(dir, "types")
	if _, err := os.Stat(typesDir); os.IsNotExist(err) {
		err = os.MkdirAll(typesDir, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// fmt.Println("PATH ", path)
	err = generator.GenerateBaseAPI(pkg, path, dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = generator.GenerateBaseTypes(pkg, path, typesDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, arg := range args {
		var content map[string]map[string]interface{}
		err = c.DoRequest("fields_get", arg, []interface{}{},
			map[string][]string{"attributes": []string{"type"}}, &content)
		if err != nil {
			if err.Error() != "error: \"\" code: 2" {
				fmt.Println(err.Error())
				return
			}
		}
		if len(content) == 0 {
			fmt.Println(fmt.Sprintf("WARN: The model %s was not found", arg))
			continue
		}
		// fmt.Println(arg, "FIELDS", content)
		err = generator.GenerateTypes(pkg, path, typesDir, arg, generateContent(content))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = generator.GenerateAPI(pkg, path, dir, arg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

func generateContent(apiContent map[string]map[string]interface{}) map[string]string {
	content := make(map[string]string, len(apiContent))
	for modelName, field := range apiContent {
		for fieldType, fieldContent := range field {
			if fieldType == "type" {
				content[modelName] = fieldContent.(string)
			}
		}
	}
	return content
}
