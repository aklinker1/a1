package cache

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
	utils "github.com/aklinker1/a1/pkg/utils"
)

// CreateDataLoader -
func CreateDataLoader(localData map[string]map[interface{}]map[string]interface{}, password string) a1.DataLoader {
	return a1.DataLoader{
		Connect: func() error {
			if password != "password" {
				return fmt.Errorf("Password authentication failed")
			}
			return nil
		},

		GetOne: func(model *a1.FinalModel, args a1.DataMap, fields a1.RequestedFieldMap) (a1.DataMap, error) {
			utils.Log("Selecting one from '%s', where %v", model.DataLoader.Group, args)
			group := localData[model.DataLoader.Group]
			var data a1.DataMap
			for _, item := range group {
				isMatch := true
				for arg, argValue := range args {
					if item[arg] != argValue {
						isMatch = false
					}
				}
				if isMatch {
					data = item
				}
			}
			utils.Log("  - Result: %v", data)
			return data, nil
		},

		GetMultiple: func(model *a1.FinalModel, args a1.DataMap, fields a1.RequestedFieldMap) ([]a1.DataMap, error) {
			utils.Log("Selecting multiple from '%s' where %v", model.DataLoader.Group, args)
			group := localData[model.DataLoader.Group]
			items := []a1.DataMap{}
			for _, item := range group {
				isMatch := true
				for arg, argValue := range args {
					if item[arg] != argValue {
						isMatch = false
					}
				}
				if isMatch {
					items = append(items, item)
				}
			}
			utils.Log("  - Result: %v", items)
			return items, nil
		},
	}
}
