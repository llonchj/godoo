package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"reflect"

	"github.com/kolo/xmlrpc"
	"github.com/llonchj/godoo/internal/types"
)

var (
	//ErrUnauthorized is a authentication error
	ErrUnauthorized = errors.New("unauthorized")
)

// CommonClient is the xmlrpc endpoint for common service
type CommonClient struct {
	*xmlrpc.Client
}

// DBClient is the xmlrpc endpoint for db service
type DBClient struct {
	*xmlrpc.Client
}

// ReportClient is the xmlrpc endpoint for report service
type ReportClient struct {
	*xmlrpc.Client
}

// Version server response
type Version struct {
	ServerVersion     string        `xmlrpc:"server_version"`
	ServerVersionInfo []interface{} `xmlrpc:"server_version_info"`
	ServerSerie       string        `xmlrpc:"server_serie"`
	ProtocolVersion   int64         `xmlrpc:"protocol_version"`
}

// Languages server response
type Languages [][]string

// Client instance for ODOO
type Client struct {
	CommonClient *CommonClient
	ObjectClient *xmlrpc.Client
	DbClient     *DBClient
	ReportClient *ReportClient

	URL *url.URL

	version *Version
}

type Session struct {
	Client   *Client
	DbName   string
	User     string
	Password string

	UID int64
}

// NewClient creates a new ODOO Client
func NewClient(URL string, Transport http.RoundTripper) (*Client, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	if Transport == nil {
		Transport = http.DefaultTransport
	}

	commonClient, err := GetCommonClient(u.String(), Transport)
	if err != nil {
		return nil, err
	}
	objectClient, err := GetObjectClient(u.String(), Transport)
	if err != nil {
		return nil, err
	}
	dbClient, err := GetDbClient(u.String(), Transport)
	if err != nil {
		return nil, err
	}
	reportClient, err := GetReportClient(u.String(), Transport)
	if err != nil {
		return nil, err
	}

	return &Client{
		CommonClient: commonClient,
		ObjectClient: objectClient,
		DbClient:     dbClient,
		ReportClient: reportClient,
		URL:          u,
	}, err
}

// NewSession creates a new ODOO Session
func (client *Client) NewSession(DbName, User, Password string) *Session {
	return &Session{
		Client:   client,
		User:     User,
		Password: Password,
		DbName:   DbName,
	}
}

// Version returns the version of the server
func (c *Client) Version() (*Version, error) {
	if c.version != nil {
		return c.version, nil
	}
	var err error
	c.version, err = c.CommonClient.Version()
	if err != nil {
		return nil, err
	}
	return c.version, err
}

// CompleteSession authenticates credentials and sets session UID
func (s *Session) CompleteSession() error {
	uid, err := s.Client.CommonClient.Authenticate(s.DbName, s.User, s.Password, nil)
	if err != nil {
		return err
	}
	s.UID = uid
	return nil
}

// GetObjectClient invokes object endpoint
func GetObjectClient(uri string, transport http.RoundTripper) (*xmlrpc.Client, error) {
	return xmlrpc.NewClient(uri+"/xmlrpc/2/object", transport)
}

// GetCommonClient invokes common endpoint
func GetCommonClient(uri string, transport http.RoundTripper) (*CommonClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/common", transport)
	if err != nil {
		return nil, err
	}
	return &CommonClient{Client: x}, nil
}

// GetDbClient invokes db endpoint
func GetDbClient(uri string, transport http.RoundTripper) (*DBClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/db", transport)
	if err != nil {
		return nil, err
	}
	return &DBClient{Client: x}, nil
}

// GetReportClient invokes report endpoint
func GetReportClient(uri string, transport http.RoundTripper) (*ReportClient, error) {
	x, err := xmlrpc.NewClient(uri+"/xmlrpc/2/report", transport)
	if err != nil {
		return nil, err
	}
	return &ReportClient{Client: x}, nil
}

// Create instantiates a new model with given args returning elem
func (c *Session) Create(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("create", model, args, options, elem)
}

// Update updates model instances
func (c *Session) Update(model string, args []interface{}, options interface{}) error {
	return c.DoRequest("write", model, args, options, nil)
}

// Delete deletes the model
func (c *Session) Delete(model string, args []interface{}, options interface{}) error {
	return c.DoRequest("unlink", model, args, options, nil)
}

// Search searches
func (c *Session) Search(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("search", model, args, options, elem)
}

// Read a element
func (c *Session) Read(model string, args []interface{}, options interface{}, elem interface{}) error {
	ne := elem.(types.Type).NilableType()
	err := c.DoRequest("read", model, args, options, ne)
	if err == nil {
		reflect.ValueOf(elem).Elem().Set(reflect.ValueOf(ne.(types.NilableType).GetType()).Elem())
	}
	return err
}

// SearchRead a
func (c *Session) SearchRead(model string, args []interface{}, options interface{}, elem interface{}) error {
	ne := elem.(types.Type).NilableType()
	err := c.DoRequest("search_read", model, args, options, ne)
	if err == nil {
		reflect.ValueOf(elem).Elem().Set(reflect.ValueOf(ne.(types.NilableType).GetType()).Elem())
	}
	return err
}

// SearchCount returns the count of records matching a search for a specified model.
func (c *Session) SearchCount(model string, args []interface{}, options interface{}, elem interface{}) error {
	return c.DoRequest("search_count", model, args, options, elem)
}

// DoRequest invokes a method over a model
func (s *Session) DoRequest(method string, model string, args []interface{}, options interface{}, elem interface{}) error {
	return s.Client.ObjectClient.Call("execute_kw",
		[]interface{}{
			s.DbName,
			s.UID,
			s.Password,
			model,
			method,
			args,
			options}, elem)
}

func (c *Session) getIdsByName(model string, name string, options interface{}) ([]int64, error) {
	var ids []int64
	err := c.Search(model, []interface{}{[]string{"name", "=", name}}, options, &ids)
	return ids, err
}

func (c *Session) getByIds(model string, ids []int64, options interface{}, elem interface{}) error {
	err := c.Read(model, []interface{}{ids}, options, elem)
	return err
}

func (c *Session) getByName(model string, name string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{[]string{"name", "=", name}}}, options, elem)
	return err
}

func (c *Session) getByField(model string, field string, value string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{[]string{field, "=", value}}}, options, elem)
	return err
}

func (c *Session) getAll(model string, options interface{}, elem interface{}) error {
	err := c.SearchRead(model, []interface{}{[]interface{}{}}, options, elem)
	return err
}

func (c *Session) create(model string, fields map[string]interface{}, relations *types.Relations, options interface{}) (int64, error) {
	var id int64
	if relations != nil {
		relations.Handle(&fields)
	}
	err := c.Create(model, []interface{}{fields}, options, &id)
	return id, err
}

func (c *Session) update(model string, ids []int64, fields map[string]interface{}, relations *types.Relations, options interface{}) error {
	if relations != nil {
		relations.Handle(&fields)
	}
	err := c.Update(model, []interface{}{ids, fields}, options)
	return err
}

func (c *Session) delete(model string, ids []int64, options interface{}) error {
	return c.Delete(model, []interface{}{ids}, options)
}

// GetAllModels returns available models
func (c *Session) GetAllModels() ([]string, error) {
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
func (c *Session) Ref(module, name string) (string, int64, error) {
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

// Authenticate authenticates credentials and sets session UID
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

// Version returns the odoo Version
func (c *CommonClient) Version() (*Version, error) {
	var version Version
	err := c.Call("version", []interface{}{}, &version)
	if err != nil {
		return nil, err
	}
	return &version, err
}

// Create invokes a method over a model
func (c *DBClient) Create(adminPassword, db string, demo bool, lang, password string, result interface{}) error {
	return c.Call("create_database",
		[]interface{}{adminPassword, db, demo, lang, password, "admin", "ES"}, result)
}

// Duplicate invokes a method over a model
func (c *DBClient) Duplicate(adminPassword, dst, src string, result interface{}) error {
	return c.Call("duplicate_database",
		[]interface{}{adminPassword, src, dst}, result)
}

// Rename invokes a method over a model
func (c *DBClient) Rename(adminPassword, dst, src string) (bool, error) {
	var result bool
	err := c.Call("rename",
		[]interface{}{adminPassword, src, dst}, &result)
	return result, err
}

// Dump invokes a method over a model
func (c *DBClient) Dump(adminPassword, db, format string) ([]byte, error) {
	var result string
	err := c.Call("dump", []interface{}{adminPassword, db, format}, &result)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(result)
	return decoded, err
}

// Restore invokes a method over a model
func (c *DBClient) Restore(adminPassword, db string, data []byte, copy bool) (bool, error) {
	var result bool
	err := c.Call("restore", []interface{}{adminPassword, db,
		base64.StdEncoding.EncodeToString(data), copy}, &result)
	return result, err
}

// List lists databases
func (c *DBClient) List() ([]string, error) {
	var dbList []string
	err := c.Call("list", []interface{}{}, &dbList)
	return dbList, err
}

// ListLanguages lists databases
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

// Exist lists databases
func (c *DBClient) Exist(db string) (bool, error) {
	var result bool
	err := c.Call("db_exist", []interface{}{db}, &result)
	return result, err
}

// Drop lists databases
func (c *DBClient) Drop(adminPassword, db string) (bool, error) {
	var result bool
	err := c.Call("drop", []interface{}{adminPassword, db}, &result)
	return result, err
}

// ServerVersion returns ODOO version
func (c *DBClient) ServerVersion() (string, error) {
	var result string
	err := c.Call("server_version", []interface{}{}, &result)
	return result, err
}

// ChangeAdminPassword changes SuperAdmin password
func (c *DBClient) ChangeAdminPassword(old, new string) (bool, error) {
	var result bool
	err := c.Call("change_admin_password", []interface{}{old, new}, &result)
	return result, err
}

// RenderReport renders a report
func (c *ReportClient) RenderReport(ReportName string, IDs []int64) ([]byte, error) {
	var result string
	err := c.Call("render_report", []interface{}{ReportName, IDs}, &result)
	decoded, err := base64.StdEncoding.DecodeString(result)
	return decoded, err
}
