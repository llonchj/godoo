package api_test

import (
	// "io/ioutil"

	"testing"
	// "github.com/seborama/govcr"
)

func TestServerVersion(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	version, err := c.DbClient.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Version", version)
}
