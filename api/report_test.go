package api_test

import (
	// "io/ioutil"

	"testing"
)

func TestReport(t *testing.T) {
	c, err := getClient(t)
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
