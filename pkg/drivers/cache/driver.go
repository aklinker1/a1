package cache

import (
	a1 "github.com/aklinker1/a1/pkg"
	utils "github.com/aklinker1/a1/pkg/utils"
)

// CreateDriver -
func CreateDriver(localData map[string]map[interface{}]map[string]interface{}) a1.DatabaseDriver {
	return a1.DatabaseDriver{
		Name:    "Cache Map",
		Connect: func() {},
		SelectOne: func(model a1.Model, primaryKey interface{}, fields a1.StringMap) (a1.StringMap, error) {
			utils.Log("Selecting one from '%s', where %s=%v (%T)", model.Table, model.PrimaryKey, primaryKey, primaryKey)
			data := localData[model.Table][primaryKey]
			utils.Log("  - Result: %v", data)
			return data, nil
		},
	}
}
