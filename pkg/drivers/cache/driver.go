package cache

import (
	a1 "github.com/aklinker1/a1/pkg/new"
	utils "github.com/aklinker1/a1/pkg/utils"
)

// CreateDriver -
func CreateDriver(localData map[string]map[interface{}]map[string]interface{}) a1.DataLoader {
	return a1.DataLoader{
		Connect: func() error {
			return nil
		},

		GetOne: func(model a1.FinalModel, primaryKey interface{}, fields a1.RequestedFieldMap) (a1.DataMap, error) {
			utils.Log("Selecting one from '%s', where %s=%v (%T)", model.DataLoader.Group, model.PrimaryKey, primaryKey, primaryKey)
			data := localData[model.DataLoader.Group][primaryKey]
			utils.Log("  - Result: %v", data)
			return data, nil
		},

		GetMultiple: func(model a1.FinalModel, searchArgs a1.DataMap, fields a1.RequestedFieldMap) ([]a1.DataMap, error) {
			utils.Log("Selecting multiple from '%s'", model.DataLoader.Group)
			itemMap := localData[model.DataLoader.Group]
			items := []a1.DataMap{}
			for _, item := range itemMap {
				items = append(items, item)
			}
			utils.Log("  - Result: %v", items)
			return items, nil
		},
	}
}
