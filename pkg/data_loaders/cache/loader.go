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

func update(model *a1.FinalModel, inputData a1.DataMap, whereArgs a1.DataMap, requestedFields a1.RequestedFieldMap) (a1.DataMap, error) {
	utils.Log("Updating one from '%s' to %v, where %v", model.DataLoader.Group, inputData, whereArgs)
	item, err := getOne(model, whereArgs, requestedFields)
	if err != nil {
		return nil, err
	}

	// Separate Linked vs Regular fields
	linkedInputFields := a1.DataMap{}
	regularInputFields := a1.DataMap{}
	for inputField, inputValue := range inputData {
		if value, isMap := inputValue.(a1.DataMap); isMap {
			linkedInputFields[inputField] = value
		} else {
			regularInputFields[inputField] = inputValue
		}
	}

	// Update the regular fields
	for fieldName, fieldValue := range regularInputFields {
		item[fieldName] = fieldValue
	}
	utils.Log("  - Result: %v", item)

	// Update the linked fields
	for fieldName, fieldValue := range linkedInputFields {
		linkedField := model.Fields[fieldName].(*a1.FinalLinkedField)
		nextModel := linkedField.LinkedModel
		nextInputData := fieldValue.(a1.DataMap)
		convertedSourceField := model.Fields[linkedField.Field].(*a1.FinalField).DataField
		convertedLinkedField := nextModel.Fields[linkedField.LinkedField].(*a1.FinalField).DataField
		nextWhereArgs := a1.DataMap{
			convertedLinkedField: item[convertedSourceField],
		}
		nextRequestedFields := requestedFields[fieldName].InnerFields.(a1.RequestedFieldMap)
		value, err := update(nextModel, nextInputData, nextWhereArgs, nextRequestedFields)
		if err != nil {
			utils.LogRed("Failed to update "+fieldName, err)
		} else {
			item[fieldName] = value
		}
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
