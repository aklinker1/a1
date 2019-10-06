package cache

import (
	"fmt"

	a1 "github.com/aklinker1/a1/pkg"
	utils "github.com/aklinker1/a1/pkg/utils"
)

var localData map[string]map[interface{}]map[string]interface{}

func getOne(model *a1.FinalModel, args a1.DataMap, fields a1.RequestedFieldMap) (a1.DataMap, error) {
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
	if len(data) == 0 {
		return nil, fmt.Errorf("No %s was found matching %v", model.Name, args)
	}
	utils.Log("  - Result: %v", data)
	return data, nil
}

func getMultipile(model *a1.FinalModel, args a1.DataMap, fields a1.RequestedFieldMap) ([]a1.DataMap, error) {
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
}

func update(model *a1.FinalModel, args a1.DataMap, inputData a1.DataMap, fields a1.RequestedFieldMap) (a1.DataMap, error) {
	item, err := getOne(model, args, fields)
	if err != nil {
		return nil, err
	}
	for inputField, inputValue := range inputData {
		item[inputField] = inputValue
	}
	return item, nil
}

// CreateDataLoader -
func CreateDataLoader(data map[string]map[interface{}]map[string]interface{}, password string) a1.DataLoader {
	return a1.DataLoader{
		Connect: func() error {
			if password != "password" {
				return fmt.Errorf("Password authentication failed")
			}
			localData = data
			return nil
		},

		GetOne:      getOne,
		GetMultiple: getMultipile,
		Update:      update,
	}
}
