# godoo

Odoo API client for Go.

## Install

```bash
go install github.com/llonchj/godoo@latest
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
- Version support (8-9-10-11)
  - Printing based on version (<https://github.com/OCA/odoorpc/blob/master/odoorpc/report.py>)
- New Odoo API functions
  - [Workflow Signals](https://www.odoo.com/documentation/8.0/api_integration.html#workflow-manipulations)
  - [Report printing](https://www.odoo.com/documentation/8.0/api_integration.html#report-printing)

```
docker compose run odoo odoo -i base -d odoo --stop-after-init
```
