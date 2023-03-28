package commands

import (
	api "github.com/llonchj/godoo/internal/api"
	"github.com/spf13/cobra"
)

func getClient(cmd *cobra.Command) (*api.Client, error) {
	uri, err := cmd.PersistentFlags().GetString("uri")
	if err != nil {
		return nil, err
	}
	return api.NewClient(uri, nil)
}

func getSession(cmd *cobra.Command) (*api.Session, error) {
	db, err := cmd.PersistentFlags().GetString("database")
	if err != nil {
		return nil, err
	}
	user, err := cmd.PersistentFlags().GetString("username")
	if err != nil {
		return nil, err
	}
	password, err := cmd.PersistentFlags().GetString("password")
	if err != nil {
		return nil, err
	}

	client, err := getClient(cmd)
	if err != nil {
		return nil, err
	}
	return client.NewSession(db, user, password), nil
}
