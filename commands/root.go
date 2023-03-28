package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godoo COMMAND [ARGS]",
	Short: "Prepare your custom odoo api environment",
	Long:  "Tool to prepare environment for odoo golang api wrapper",
}

// Execute is the entry point for the command line
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
