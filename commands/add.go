package commands

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [args]",
	Short: "add the given odoo model",
	Long:  "This command add the needed api and type files for the given odoo model",
	Example: `
	./godoo add account.analytic.account
	./godoo add account.analytic.account account.invoice
	./godoo add all`,
	Run: upsert,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringP("uri", "", "http://localhost:8069", "the odoo instance URI")
	addCmd.PersistentFlags().StringP("database", "d", "database", "the name of the postgres database linked to odoo")
	addCmd.PersistentFlags().StringP("username", "u", "admin", "the odoo instance administrator")
	addCmd.PersistentFlags().StringP("password", "p", "password", "the odoo instance administrator password")
	addCmd.PersistentFlags().StringP("package", "", "godoo", "go package name")
	addCmd.PersistentFlags().StringP("path", "", "", "go package to generate the files")
}
