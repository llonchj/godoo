//go:generate godoo add all --package godoo_project --uri http://localhost:8069 -d odoo -u admin -p admin

package main

import (
	"fmt"

	"github.com/llonchj/godoo/test_project/godoo_project"
)

func main() {
	// create a Client instance
	client, err := godoo_project.NewClient("http://localhost:8069", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	// create a Session instance
	session := &godoo_project.Session{
		DbName:   "odoo",
		User:     "admin",
		Password: "admin",
		Client:   client,
	}
	err = session.CompleteSession()
	if err != nil {
		fmt.Println(err.Error())
	}

	//get the res.users service
	s := godoo_project.NewResUsersService(session)
	//call the function GetAll() linked to the sale order service
	so, err := s.GetAll(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(so)
}
