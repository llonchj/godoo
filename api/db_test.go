package api_test

import (
	// "io/ioutil"

	"archive/zip"
	"bytes"
	"testing"
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

func TestChangeAdminPassword(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	result, err := c.DbClient.ChangeAdminPassword(SuperAdminPassword, SuperAdminPassword)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Change Admin Password", result)

	_, err = c.DbClient.ChangeAdminPassword("invalid-admin", SuperAdminPassword)
	if err == nil {
		t.Fatalf("should fail here")
	}
	t.Logf("Change Admin Password %+v", err.Error())
}

func TestListDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	dbList, err := c.DbClient.List()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Database List", dbList)
}

func TestListLanguages(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.DbClient.ListLanguages()
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Language List", result)
}

func TestCreateDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	var result interface{}
	if err := c.DbClient.Create(SuperAdminPassword, TestDBPrefix+"test", true,
		"es_ES", DefaultPassword, &result); err != nil {
		t.Log("C ", result)
		t.Fatal(err)
	}
	t.Log("C ", result)
}

func TestExistsDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	exist, err := c.DbClient.Exist(TestDBPrefix + "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("D ", exist)
}

func TestAuthenticate(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}
	s := getSession(t, c)

	UID, err := c.CommonClient.Authenticate(s.DbName, s.User, s.Password, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UID", UID)

	UID, err = c.CommonClient.Authenticate(s.DbName, s.User, "Invalid", nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDuplicateDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	var result interface{}
	if err := c.DbClient.Duplicate(SuperAdminPassword, TestDBPrefix+"test_duplicate", TestDBPrefix+"test", &result); err != nil {
		t.Fatal(err)
	}
	t.Log("D ", result)
}

func TestRenameDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	result, err := c.DbClient.Rename(SuperAdminPassword,
		TestDBPrefix+"test_renamed",
		TestDBPrefix+"test_duplicate")
	if err != nil {
		t.Log("D ", result)
		t.Fatal(err)
	}
	t.Log("D ", result)
}

func TestDumpRestoreDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	data, err := c.DbClient.Dump(SuperAdminPassword, TestDBPrefix+"test", "zip")
	if err != nil {
		t.Fatal(err)
	}

	// if err := ioutil.WriteFile("dump.zip", data, 0644); err != nil {
	// 	t.Fatal(err)
	// }
	if _, err := zip.NewReader(bytes.NewReader(data), int64(len(data))); err != nil {
		t.Fatal(err)
	}

	result, err := c.DbClient.Restore(SuperAdminPassword, TestDBPrefix+"test_restored", data, false)
	if err != nil {
		t.Log("D ", result)
		t.Fatal(err)
	}
	t.Log("D ", result)

	if _, err := c.DbClient.Drop(SuperAdminPassword, TestDBPrefix+"test_restored"); err != nil {
		t.Fatal(err)
	}
}

func TestDropDB(t *testing.T) {
	c, err := getClient(t)
	if err != nil {
		t.Fatal(err)
	}

	for _, dbname := range []string{"test_renamed", "test"} {
		result, err := c.DbClient.Drop(SuperAdminPassword, TestDBPrefix+dbname)
		if err != nil {
			t.Log("D ", result)
			t.Fatal(err)
		}
		t.Log("D ", result)
	}
}
