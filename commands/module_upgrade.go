package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var moduleUpgradeCmd = &cobra.Command{
	Use:   "module-upgrade [modules]",
	Short: "upgrades specified modules",
	Long:  "This command upgrades specified modules into odoo",
	Example: `
	./godoo module-upgrade stock
	./godoo module-upgrade stock mail`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("At least one argument needed")
		}
		return nil
	},
	Run: immediateUpgrade,
}

func init() {
	rootCmd.AddCommand(moduleUpgradeCmd)
	moduleUpgradeCmd.PersistentFlags().StringP("uri", "", "http://localhost:8069", "the odoo instance URI")
	moduleUpgradeCmd.PersistentFlags().StringP("database", "d", "database", "the name of the postgres database linked to odoo")
	moduleUpgradeCmd.PersistentFlags().StringP("username", "u", "admin", "the odoo instance administrator")
	moduleUpgradeCmd.PersistentFlags().StringP("password", "p", "password", "the odoo instance administrator password")
}

func immediateUpgrade(cmd *cobra.Command, args []string) {
	moduleAction("button_immediate_upgrade", cmd, args)
}
