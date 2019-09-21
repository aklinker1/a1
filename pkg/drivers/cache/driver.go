package cache

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
)

// CreateDriver -
func CreateDriver(localData map[string]map[interface{}]interface{}) a1.DatabaseDriver {
	return a1.DatabaseDriver{
		Name:    "Cache Map",
		Connect: func() {},
		SelectOne: func(model a1.Model, id interface{}, fields a1.StringMap) (a1.StringMap, error) {
			output := a1.StringMap{}
			todo, ok := localData[model.Name][id]
			if !ok {
				return nil, fmt.Errorf("'%s' with %s=%v does not exist", model.Name, model.PrimaryKey, id)
			}
			output["id"] = id
			output["title"] = todo
			return output, nil
		},
	}
}
