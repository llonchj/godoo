package api

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/kolo/xmlrpc"
	"github.com/llonchj/godoo/types"
)

//Config structure for Odoo
type Config struct {
	DbName    string
	User      string
	Password  string
	URI       string
	Transport http.RoundTripper
}

//Client instance for ODOO
type Client struct {
	CommonClient *xmlrpc.Client
	ObjectClient *xmlrpc.Client
	DbClient     *xmlrpc.Client
	ReportClient *xmlrpc.Client
	Session      *Session
}

//Session instance for ODOO
type Session struct {
	DbName   string
	User     string
	Password string
	UID      int64
}

//NewClient creates a new ODOO Client
func (config *Config) NewClient() (*Client, error) {
	commonClient, err := GetCommonClient(config.URI, config.Transport)
	if err != nil {
		return nil, err
	}
	objectClient, err := GetObjectClient(config.URI, config.Transport)
	if err != nil {
		return nil, err
	}
	dbClient, err := GetDbClient(config.URI, config.Transport)
	if err != nil {
		return nil, err
	}
	reportClient, err := GetObjectClient(config.URI, config.Transport)
	if err != nil {
		return nil, err
	}
	return &Client{
		CommonClient: commonClient,
		ObjectClient: objectClient,
		DbClient:     dbClient,
		ReportClient: reportClient,
		Session: &Session{
			User:     config.User,
			Password: config.Password,
			DbName:   config.DbName,
		},
	}, err
}

//CompleteSession authenticates credentials and sets session UID
func (c *Client) CompleteSession() error {
	var uid interface{}
	var ok bool
	err := c.CommonClient.Call("authenticate", []interface{}{c.Session.DbName, c.Session.User, c.Session.Password, ""}, &uid)
	if err != nil {
		return err
	}
	i, ok := uid.(int64)
	if !ok {
		return errors.New("invalid session uid")
	}
	c.Session.UID = i
	return nil
}

//GetObjectClient invokes object endpoint
func GetObjectClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/object", transport)
}

//GetCommonClient invokes common endpoint
func GetCommonClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/common", transport)
}

//GetDbClient invokes db endpoint
func GetDbClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/db", transport)
}

//GetReportClient invokes report endpoint
func GetReportClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/report", transport)
}

//Create instantiates a new model with given args returning elem
func (c *Client) Create(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("create", model, args, options, elem)
}

//Update updates model instances
func (c *Client) Update(model string, args []interface{}, options interface{}) error {
	return c.DoRequest("write", model, args, options, nil)
}

//Delete deletes the model
func (c *Client) Delete(model string, args []interface{}, options interface{}) error {
	return c.DoRequest("unlink", model, args, options, nil)
}

//Search searches
func (c *Client) Search(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("search", model, args, options, elem)
}

//Read a element
func (c *Client) Read(model string, args []interface{}, options interface{}, elem interface{}) error {
	ne := elem.(types.Type).NilableType()
	err := c.DoRequest("read", model, args, options, ne)
	if err == nil {
		reflect.ValueOf(elem).Elem().Set(reflect.ValueOf(ne.(types.NilableType).GetType()).Elem())
	}
	return err
}

//SearchRead a
func (c *Client) SearchRead(model string, args []interface{}, options interface{}, elem interface{}) error {
	ne := elem.(types.Type).NilableType()
	err := c.DoRequest("search_read", model, args, options, ne)
	if err == nil {
		reflect.ValueOf(elem).Elem().Set(reflect.ValueOf(ne.(types.NilableType).GetType()).Elem())
	}
	return err
}

//SearchCount returns the count of records matching a search for a specified model.
func (c *Client) SearchCount(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("search_count", model, args, options, elem)
}

//DoRequest invokes a method over a model
func (c *Client) DoRequest(method string, model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.ObjectClient.Call("execute_kw",
		[]interface{}{
			c.Session.DbName,
			c.Session.UID,
			c.Session.Password,
			model,
			method,
			args,
			options}, elem)
}

func (c *Client) getIdsByName(model string, name string, options interface{}) ([]int64, error) {
	var ids []int64
	err := c.Search(model, []interface{}{[]string{"name", "=", name}}, options, &ids)
	return ids, err
}

func (c *Client) getByIds(model string, ids []int64, options interface{}, elem interface{}) error {
	err := c.Read(model, []interface{}{ids}, options, elem)
	return err
}

func (c *Client) getByName(model string, name string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{[]string{"name", "=", name}}}, options, elem)
	return err
}

func (c *Client) getByField(model string, field string, value string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{[]string{field, "=", value}}}, options, elem)
	return err
}

func (c *Client) getAll(model string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{}}, options, elem)
	return err
}

func (c *Client) create(model string, fields map[string]interface{}, relation *types.Relations, options interface{}) (int64, error) {
	var id int64
	if relation != nil {
		types.HandleRelations(&fields, relation)
	}
	err := c.Create(model, []interface{}{fields}, options, &id)
	return id, err
}

func (c *Client) update(model string, ids []int64, fields map[string]interface{}, relation *types.Relations, options interface{}) error {
	if relation != nil {
		types.HandleRelations(&fields, relation)
	}
	err := c.Update(model, []interface{}{ids, fields}, options)
	return err
}

func (c *Client) delete(model string, ids []int64, options interface{}) error {
	return c.Delete(model, []interface{}{ids}, options)
}

//GetAllModels returns available models
func (c *Client) GetAllModels() ([]string, error) {
	var content []map[string]interface{}
	err := c.DoRequest("search_read", "ir.model", []interface{}{[]interface{}{}}, nil, &content)
	if err != nil {
		return []string{}, err
	}
	models := make([]string, len(content))
	for i, modelFields := range content {
		for field, model := range modelFields {
			if field == "model" {
				models[i] = model.(string)
			}
		}
	}
	return models, err
}

// Ref returns the model,id of a external identifier. This function is inspired
// in Odoo new api self.env.ref function.
func (c *Client) Ref(module, name string) (string, int64, error) {
	var content []map[string]interface{}
	err := c.DoRequest("search_read", "ir.model.data",
		[]interface{}{
			[]interface{}{
				[]string{"module", "=", module},
				[]string{"name", "=", name},
			},
		}, nil, &content)
	if err != nil {
		return "", 0, err
	}
	return content[0]["model"].(string), content[0]["res_id"].(int64), nil
}
