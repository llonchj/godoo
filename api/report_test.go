package api_test

import (
	// "io/ioutil"

	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/llonchj/godoo/api"
	"github.com/seborama/govcr"
)

func TestReport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			fmt.Fprintln(w, `<?xml version='1.0'?><methodResponse>
			<params><param><value><string>TEST</string></value></param></params>
			</methodResponse>
			`)
		}))
	defer ts.Close()

	// t.Logf("odoo version does not support reports")

	vcr := govcr.NewVCR(t.Name(),
		&govcr.VCRConfig{
			// Logging:          true,
			// DisableRecording: false,
			CassettePath: "./fixtures",
		})

	config := api.Config{
		DbName: "test",
		// URI:      ts.URL, //"http://localhost:8069",
		URI:       "http://localhost:8069",
		User:      "admin",
		Password:  "password",
		Transport: vcr.Client.Transport,
	}
	c, err := config.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	reportName := "base.report_irmodeloverview"
	ids := []int64{1, 2, 3}

	_, err = c.ReportClient.RenderReport(reportName, ids)
	if err != nil {
		t.Fatal(err)
	}
	// ioutil.WriteFile("data.zip", data, 0644)
}
