package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"reflect"

	"github.com/kolo/xmlrpc"
	"github.com/llonchj/godoo/types"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

//Config structure for Odoo
type Config struct {
	DbName    string
	User      string
	Password  string
	URI       string
	Transport http.RoundTripper
}

//CommonClient
type CommonClient struct {
	*xmlrpc.Client
}

//DBClient
type DBClient struct {
	*xmlrpc.Client
}

//ReportClient
type ReportClient struct {
	*xmlrpc.Client
}

//Client instance for ODOO
type Client struct {
	CommonClient *CommonClient
	ObjectClient *xmlrpc.Client
	DbClient     *DBClient
	ReportClient *ReportClient
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
	reportClient, err := GetReportClient(config.URI, config.Transport)
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
	uid, err := c.CommonClient.Authenticate(c.Session.DbName, c.Session.User, c.Session.Password, nil)
	if err != nil {
		return err
	}
	c.Session.UID = uid
	return nil
}

//GetObjectClient invokes object endpoint
func GetObjectClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/object", transport)
}

//GetCommonClient invokes common endpoint
func GetCommonClient(uri string, transport http.RoundTripper) (*CommonClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/common", transport)
	if err != nil {
		return nil, err
	}
	return &CommonClient{Client: x}, nil
}

//GetDbClient invokes db endpoint
func GetDbClient(uri string, transport http.RoundTripper) (*DBClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/db", transport)
	if err != nil {
		return nil, err
	}
	return &DBClient{Client: x}, nil
}

//GetReportClient invokes report endpoint
func GetReportClient(uri string, transport http.RoundTripper) (*ReportClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/report", transport)
	if err != nil {
		return nil, err
	}
	return &ReportClient{Client: x}, nil
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

type Version struct {
	ServerVersion     string        `xmlrpc:"server_version"`
	ServerVersionInfo []interface{} `xmlrpc:"server_version_info"`
	ServerSerie       string        `xmlrpc:"server_serie"`
	ProtocolVersion   int64         `xmlrpc:"protocol_version"`
}

//Authenticate
//CompleteSession authenticates credentials and sets session UID
func (c *CommonClient) Authenticate(DbName, User, Password string, UserAgentEnv map[string]interface{}) (int64, error) {
	var uid interface{}
	err := c.Call("authenticate", []interface{}{DbName, User, Password, UserAgentEnv}, &uid)
	if err != nil {
		return 0, err
	}

	i, ok := uid.(int64)
	if !ok {
		//Assume Odoo returned false
		return 0, ErrUnauthorized
	}
	return i, nil
}

//Version
func (c *CommonClient) Version() (*Version, error) {
	var version Version
	err := c.Call("version", []interface{}{}, &version)
	if err != nil {
		return nil, err
	}
	return &version, err
}

//Create invokes a method over a model
func (c *DBClient) Create(adminPassword, db string, demo bool, lang, password string, result interface{}) error {
	return c.Call("create_database",
		[]interface{}{adminPassword, db, demo, lang, password, "admin", "ES"}, result)
}

//Duplicate invokes a method over a model
func (c *DBClient) Duplicate(adminPassword, dst, src string, result interface{}) error {
	return c.Call("duplicate_database",
		[]interface{}{adminPassword, src, dst}, result)
}

//Rename invokes a method over a model
func (c *DBClient) Rename(adminPassword, dst, src string) (bool, error) {
	var result bool
	err := c.Call("rename",
		[]interface{}{adminPassword, src, dst}, &result)
	return result, err
}

//Dump invokes a method over a model
func (c *DBClient) Dump(adminPassword, db, format string) ([]byte, error) {
	var result string
	err := c.Call("dump", []interface{}{adminPassword, db, format}, &result)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(result)
	return decoded, err
}

//Restore invokes a method over a model
func (c *DBClient) Restore(adminPassword, db string, data []byte, copy bool) (bool, error) {
	var result bool
	err := c.Call("restore", []interface{}{adminPassword, db,
		base64.StdEncoding.EncodeToString(data), copy}, &result)
	return result, err
}

//List lists databases
func (c *DBClient) List() ([]string, error) {
	var dbList []string
	err := c.Call("list", []interface{}{}, &dbList)
	return dbList, err
}

type Languages [][]string

//ListLanguages lists databases
func (c *DBClient) ListLanguages() (Languages, error) {
	var languages Languages
	err := c.Call("list_lang", []interface{}{}, &languages)
	return languages, err
}

// type Countries interface{}

// //ListCountries lists databases
// func (c *Client) ListCountries() (Countries, error) {
// 	var countries Countries
// 	err:= c.DbClient.Call("list_countries", []interface{}{}, &countries)
// 	return countries, err
// }

//Exist lists databases
func (c *DBClient) Exist(db string) (bool, error) {
	var result bool
	err := c.Call("db_exist", []interface{}{db}, &result)
	return result, err
}

//Drop lists databases
func (c *DBClient) Drop(adminPassword, db string) (bool, error) {
	var result bool
	err := c.Call("drop", []interface{}{adminPassword, db}, &result)
	return result, err
}

//ServerVersion returns ODOO version
func (c *DBClient) ServerVersion() (string, error) {
	var result string
	err := c.Call("server_version", []interface{}{}, &result)
	return result, err
}

//ChangeAdminPassword changes SuperAdmin password
func (c *DBClient) ChangeAdminPassword(old, new string) (bool, error) {
	var result bool
	err := c.Call("change_admin_password", []interface{}{old, new}, &result)
	return result, err
}

//RenderReport renders a report
func (c *ReportClient) RenderReport(ReportName string, IDs []int64) ([]byte, error) {
	var result string
	err := c.Call("render_report", []interface{}{ReportName, IDs}, &result)
	decoded, err := base64.StdEncoding.DecodeString(result)
	return decoded, err
}
