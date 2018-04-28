package {{.Package}}

import (
	"{{.Path}}/{{.Package}}/types"
)

type {{ funcName .Model}}Service struct {
	client *Client
}

func New{{ funcName .Model}}Service(c *Client) *{{ funcName .Model}}Service {
	return &{{ funcName .Model}}Service{client: c}
}

func (svc *{{ funcName .Model}}Service) GetIdsByName(name string) ([]int64, error) {
	return svc.client.getIdsByName(types.{{ funcName .Model}}Model, name)
}

func (svc *{{ funcName .Model}}Service) GetByIds(ids []int64) (*types.{{ funcName .Model}}s, error) {
	p := &types.{{ funcName .Model}}s{}
	return p, svc.client.getByIds(types.{{ funcName .Model}}Model, ids, p)
}

func (svc *{{ funcName .Model}}Service) GetByName(name string) (*types.{{ funcName .Model}}s, error) {
	p := &types.{{ funcName .Model}}s{}
	return p, svc.client.getByName(types.{{ funcName .Model}}Model, name, p)
}

func (svc *{{ funcName .Model}}Service) GetByField(field string, value string) (*types.{{ funcName .Model}}s, error) {
	p := &types.{{ funcName .Model}}s{}
	return p, svc.client.getByField(types.{{ funcName .Model}}Model, field, value, p)
}

func (svc *{{ funcName .Model}}Service) GetAll() (*types.{{ funcName .Model}}s, error) {
	p := &types.{{ funcName .Model}}s{}
	return p, svc.client.getAll(types.{{ funcName .Model}}Model, p)
}

func (svc *{{ funcName .Model}}Service) Create(fields map[string]interface{}, relations *types.Relations) (int64, error) {
	return svc.client.create(types.{{ funcName .Model}}Model, fields, relations)
}

func (svc *{{ funcName .Model}}Service) Update(ids []int64, fields map[string]interface{}, relations *types.Relations) error {
	return svc.client.update(types.{{ funcName .Model}}Model, ids, fields, relations)
}

func (svc *{{ funcName .Model}}Service) Delete(ids []int64) error {
	return svc.client.delete(types.{{ funcName .Model}}Model, ids)
}
