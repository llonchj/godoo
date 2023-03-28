package commands

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var installLanguageCmd = &cobra.Command{
	Use:   "install-language [languages]",
	Short: "installs specified languages",
	Long:  "This command loads the specified languages into odoo",
	Example: `
	./godoo install-language en_US
	./godoo install-language en_US zh_CN`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("At least one argument needed")
		}
		return nil
	},
	Run: installLanguage,
}

func init() {
	rootCmd.AddCommand(installLanguageCmd)
	installLanguageCmd.PersistentFlags().StringP("uri", "", "http://localhost:8069", "the odoo instance URI")
	installLanguageCmd.PersistentFlags().StringP("database", "d", "database", "the name of the postgres database linked to odoo")
	installLanguageCmd.PersistentFlags().StringP("username", "u", "admin", "the odoo instance administrator")
	installLanguageCmd.PersistentFlags().StringP("password", "p", "password", "the odoo instance administrator password")
	installLanguageCmd.PersistentFlags().BoolP("overwrite", "o", false, "overwrite existing terms")
}

func installLanguage(cmd *cobra.Command, args []string) {
	c, err := getSession(cmd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := c.CompleteSession(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	command := "lang_install"
	model := "base.language.install"

	var wg sync.WaitGroup

	overwrite, err := cmd.PersistentFlags().GetBool("overwrite")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, lang := range args {
		wg.Add(1)
		go func(lang string) {
			defer wg.Done()
			var ID int64
			if err := c.Create(model, []interface{}{
				map[string]interface{}{
					"lang":      lang,
					"overwrite": overwrite,
				},
			}, nil, &ID); err != nil {
				panic(err)
			}

			var result interface{}
			if err := c.DoRequest(command, model,
				[]interface{}{[]int64{ID}},
				nil, &result); err != nil {
				panic(err)
			}

			fmt.Println(result)
		}(lang)
	}
	wg.Wait()
}
