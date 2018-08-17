package types

import (
	"reflect"
	"time"
)

type Code int64

const (
	AddRecordCode      Code = 0
	UpdateRecordCode   Code = 1
	DeleteRecordCode   Code = 2
	DeleteIDCode       Code = 3
	AddIDCode          Code = 4
	RemoveAllIDsCode   Code = 5
	ReplaceWithIDsCode Code = 6
)

//Command represents a interface for mutations
type Command interface {
	Command() []interface{}
}

type Commands []Command

//Commands returns the list of commands as Command structure suitable for odoo
//One2Many and Many2Many create/write operation
func (c *Commands) Commands() []interface{} {
	result := []interface{}{}
	for _, command := range *c {
		result = append(result, command.Command())
	}
	return result
}

type Record map[string]interface{}

type AddRecord []Record

func (r AddRecord) Command() []interface{} { return []interface{}{AddRecordCode, nil, r} }

type UpdateRecord struct {
	ID     int64
	Values Record
}

func (r UpdateRecord) Command() []interface{} {
	return []interface{}{UpdateRecordCode, r.ID, r.Values}
}

type DeleteID int64

func (r DeleteID) Command() []interface{} {
	return []interface{}{DeleteIDCode, r, nil}
}

type DeleteRecord int64

func (r DeleteRecord) Command() []interface{} {
	return []interface{}{DeleteRecordCode, r, nil}
}

type AddID int64

func (r AddID) Command() []interface{} { return []interface{}{AddIDCode, r, nil} }

type RemoveAllIDs struct{}

func (r RemoveAllIDs) Command() []interface{} {
	return []interface{}{RemoveAllIDsCode, nil, nil}
}

type ReplaceWithIDs []int64

func (r ReplaceWithIDs) Command() []interface{} {
	return []interface{}{ReplaceWithIDsCode, nil, r}
}

//Relations is a list of Relation
type Relations map[string]Commands

//NilableType is ...
type NilableType interface {
	GetType() interface{}
}

//Type is ...
type Type interface {
	NilableType() interface{}
}

//Many2One is ...
type Many2One struct {
	ID   int64
	Name string
}

func load(ns interface{}, s interface{}) interface{} {
	nse := reflect.ValueOf(ns).Elem()
	se := reflect.ValueOf(s).Elem()
	for i := 0; i < nse.NumField(); i++ {
		f := nse.Field(i)
		if f.Kind() == reflect.Bool {
			se.Field(i).SetBool(f.Bool())
			continue
		}
		if f.IsNil() || f.Elem().Kind() == reflect.Bool {
			continue
		}
		if se.Field(i).Type().Name() == "Time" {
			var t time.Time
			d := f.Elem().String()
			if len(d) == 10 {
				t, _ = time.Parse("2006-01-02", d)
			} else {
				t, _ = time.Parse("2006-01-02 15:04:05", d)
			}
			se.Field(i).Set(reflect.ValueOf(t))
			continue
		}
		if f.Elem().Kind() == reflect.Slice {
			if se.Field(i).Kind() == reflect.Slice {
				interfaceSlice := f.Interface().([]interface{})
				int64Slice := make([]int64, len(interfaceSlice))
				for j := range interfaceSlice {
					int64Slice[j] = interfaceSlice[j].(int64)
				}
				se.Field(i).Set(reflect.ValueOf(int64Slice))
				continue
			}
			if se.Field(i).Kind() == reflect.Struct {
				se.Field(i).Set(reflect.ValueOf(Many2One{ID: f.Elem().Index(0).Elem().Int(), Name: f.Elem().Index(1).Elem().String()}))
				continue
			}
		}
		se.Field(i).Set(f.Elem())
	}
	return s
}

//Handle completes fields structure with relations
func (r *Relations) Handle(fields *map[string]interface{}) {
	for field, commands := range *r {
		(*fields)[field] = commands.Commands()
	}
}
