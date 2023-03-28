package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var serverVersionCmd = &cobra.Command{
	Use:   "server-version",
	Short: "returns the version for the specified server",
	Long:  "This command queries the ODOO version for the specified server",
	Example: `
	./godoo server-version`,
	Run: serverVersion,
}

func init() {
	rootCmd.AddCommand(serverVersionCmd)
	serverVersionCmd.PersistentFlags().StringP("uri", "", "http://localhost:8069", "the odoo instance URI")
}

func serverVersion(cmd *cobra.Command, args []string) {
	c, err := getClient(cmd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	version, err := c.CommonClient.Version()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Server version:", version.ServerVersion)
	fmt.Println("Server serie:", version.ServerSerie)
	fmt.Println("Protocol version:", version.ProtocolVersion)
}
