# godoo

	go generate; and go install .

> This project is a fork of [github.com/antony360/go-odoo](https://github.com/antony360/go-odoo)

A Odoo API client enabling Go programs to interact with Odoo in a simple and uniform way.

## Coverage

This API client package covers all basic functions from the odoo API.
This include all calls to the following services :

- [x] Login
- [x] Create
- [x] Update
- [x] Delete
- [x] Search
- [x] Read
- [x] SearchRead
- [x] SearchCount
- [x] DoRequest

Services listed above are basic low-level functions from Odoo API, there accessible by any client.

There also are high-level functions based on these low-level functions. Each model has its own functions.
Actually we got:

- [x] GetIdsByName
- [x] GetByIds
- [x] GetByName
- [x] GetByField
- [x] GetAll
- [x] Create
- [x] Update
- [x] Delete

## Install

```
go get github.com/mjibson/esc github.com/llonchj/godoo
cd $GOPATH/github.com/llonchj/godoo
go generate
go install
```

## Using the API

This is an example of how create a new sale order :

```go
package main

//go:generate godoo add all --uri http://localhost:8069 --database test --username admin --password password


import (
	"fmt"

	"github.com/llonchj/godoo/api"
)

func main() {
	config := &api.Config{
		DbName:   "test",
		Admin:    "admin",
		Password: "password",
		URI:      "http://localhost:8069",
	}

	c, err := config.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.CompleteSession()
	if err != nil {
		fmt.Println(err.Error())
	}

	//get the sale order service
	s := api.NewSaleOrderService(c)
	//call the function GetAll() linked to the sale order service
	so, err := s.GetAll()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(so)
}
```

## ToDo

- Tests
- Documentation in generated files
- New Odoo API functions (ex: report printing)
- Database services
