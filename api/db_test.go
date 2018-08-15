package api_test

import (
	// "io/ioutil"
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/llonchj/godoo/api"
	// "github.com/seborama/govcr"
)

func TestServerVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/xml; charset=utf-8")
			fmt.Fprintln(w, `<?xml version='1.0'?><methodResponse>
			<params><param><value><string>TEST</string></value></param></params>
			</methodResponse>
			`)
		}))
	defer ts.Close()

	// vcr := govcr.NewVCR(t.Name(),
	// 	&govcr.VCRConfig{
	// 		Logging:          true,
	// 		DisableRecording: false,
	// 		CassettePath:     "./fixtures",
	// 	})

	config := api.Config{
		DbName:   "test",
		URI:      ts.URL, //"http://localhost:8069",
		User:     "admin",
		Password: "password",
		// Transport: vcr.Client.Transport,
	}
	c, err := config.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	version, err := c.DbClient.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}

	if version != "TEST" {
		t.Errorf("Invalid version response: got '%s'", version)
	}
	t.Log("Version", version)
}

// func TestChangeAdminPassword(t *testing.T) {
// 	config := api.Config{
// 		DbName:   "test",
// 		URI:      "http://localhost:8069",
// 		User:     "admin",
// 		Password: "password",
// 	}
// 	c, err := config.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	result, err := c.DbClient.ChangeAdminPassword("admin", "admin")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log("Change Admin Password", result)
// }

// func TestListDB(t *testing.T) {
// 	config := api.Config{
// 		DbName:   "test",
// 		URI:      "http://localhost:8069",
// 		User:     "admin",
// 		Password: "password",
// 	}
// 	c, err := config.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dbList, err := c.DbClient.List()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log("Database List", dbList)
// }

// func TestListLanguages(t *testing.T) {
// 	config := api.Config{
// 		DbName:   "test",
// 		URI:      "http://localhost:8069",
// 		User:     "admin",
// 		Password: "password",
// 	}
// 	c, err := config.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = c.DbClient.ListLanguages()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// t.Log("Language List", result)
// }

// // func TestListCountries(t *testing.T) {
// // 	config := api.Config{
// // 		DbName:   "test",
// // 		URI:      "http://localhost:8069",
// // 		User:     "admin",
// // 		Password: "password",
// // 	}
// // 	c, err := config.NewClient()
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}

// // 	result,err := c.ListCountries()
// // 	if err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	t.Log("Country List", result)
// // }

// func TestCreateDB(t *testing.T) {
// 	config := api.Config{
// 		DbName:   "test",
// 		URI:      "http://localhost:8069",
// 		User:     "admin",
// 		Password: "password",
// 	}
// 	c, err := config.NewClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// var result interface{}
// 	// if err := c.DbClient.Create("admin", "test_db", true,
// 	// 	"es_ES", "password", &result); err != nil {
// 	// 	t.Log("C ", result)
// 	// 	t.Fatal(err)
// 	// }
// 	// t.Log("C ", result)

// 	// if err := c.DbClient.Duplicate("admin", "jibu_duplicate", "jibu", &result); err != nil {
// 	// 	t.Log("D ", result)
// 	// 	t.Fatal(err)
// 	// }
// 	// t.Log("D ", result)

// 	// exist, err := c.DbClient.Exist("jibu")
// 	// if err != nil {
// 	// 	t.Log("D ", result)
// 	// 	t.Fatal(err)
// 	// }
// 	// t.Log("D ", exist)

// 	// result, err := c.DbClient.Drop("admin","jibu_duplicate")
// 	// if err != nil {
// 	// 	t.Log("D ", result)
// 	// 	t.Fatal(err)
// 	// }
// 	// t.Log("D ", result)

// 	// result, err := c.DbClient.Rename("admin","test_db2", "test_db")
// 	// if err != nil {
// 	// 	t.Log("D ", result)
// 	// 	t.Fatal(err)
// 	// }
// 	// t.Log("D ", result)

// 	data, err := c.DbClient.Dump("admin", "test_db2", "zip")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// if err:=ioutil.WriteFile("dump.zip", data, 0644);err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	result, err := c.DbClient.Restore("admin", "test_db3", data, false)
// 	if err != nil {
// 		t.Log("D ", result)
// 		t.Fatal(err)
// 	}
// 	t.Log("D ", result)

// 	// if err := c.DropDB("admin","jibu_duplic2ate"); err != nil {
// 	// 	t.Fatal(err)
// 	// }
// }
