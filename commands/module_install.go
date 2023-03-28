package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var moduleInstallCmd = &cobra.Command{
	Use:   "module-install [modules]",
	Short: "installs specified modules",
	Long:  "This command installs specified modules into odoo",
	Example: `
	./godoo module-install stock
	./godoo module-install stock mail`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("At least one argument needed")
		}
		return nil
	},
	Run: immediateInstall,
}

func init() {
	rootCmd.AddCommand(moduleInstallCmd)
	moduleInstallCmd.PersistentFlags().StringP("uri", "", "http://localhost:8069", "the odoo instance URI")
	moduleInstallCmd.PersistentFlags().StringP("database", "d", "database", "the name of the postgres database linked to odoo")
	moduleInstallCmd.PersistentFlags().StringP("username", "u", "admin", "the odoo instance administrator")
	moduleInstallCmd.PersistentFlags().StringP("password", "p", "password", "the odoo instance administrator password")
}

func immediateInstall(cmd *cobra.Command, args []string) {
	moduleAction("button_immediate_install", cmd, args)
}
