package api_test

import (
	// "io/ioutil"

	"testing"

	"github.com/llonchj/godoo/api"
	"github.com/seborama/govcr"
)

const (
	SuperAdminPassword  = "admin"
	TestDBPrefix        = "godoo_"
	DefaultUser         = "admin"
	DefaultPassword     = "password"
	DefaultURL          = "http://localhost:8069"
	VCRLogging          = true
	VCRDisableRecording = true
)

func getClient(t *testing.T) (*api.Client, error) {
	vcr := govcr.NewVCR(t.Name(), &govcr.VCRConfig{
		Logging:          VCRLogging,
		DisableRecording: VCRDisableRecording,
		CassettePath:     "./fixtures",
	})
	return api.NewClient(DefaultURL, vcr.Client.Transport)
}

func getSession(t *testing.T, client *api.Client) *api.Session {
	return client.NewSession(
		TestDBPrefix+"test", DefaultUser, DefaultPassword,
	)
}

func TestCommonVersion(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	version, err := c.CommonClient.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Version", version)
}
