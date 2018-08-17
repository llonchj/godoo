package types_test

import (
	"reflect"
	"testing"

	"github.com/llonchj/godoo/types"
)

func TestTypes(t *testing.T) {
	for _, tt := range []struct {
		Fields    map[string]interface{}
		Relations types.Relations
		Expected  map[string]interface{}
	}{
		{
			Relations: types.Relations{
				"field": types.Commands{types.AddID(0)},
			},
			Expected: map[string]interface{}{
				"field": []interface{}{[]interface{}{types.AddIDCode,
					types.AddID(0), nil}},
			},
		},
		{
			Relations: types.Relations{
				"field": types.Commands{types.RemoveAllIDs{}},
			},
			Expected: map[string]interface{}{
				"field": []interface{}{[]interface{}{types.RemoveAllIDsCode, nil, nil}},
			},
		},
	} {
		if tt.Fields == nil {
			tt.Fields = map[string]interface{}{}
		}
		tt.Relations.Handle(&tt.Fields)
		if !reflect.DeepEqual(tt.Fields, tt.Expected) {
			t.Errorf("Not equal, got %+v, expected %+v", tt.Fields, tt.Expected)
		}
	}
}
