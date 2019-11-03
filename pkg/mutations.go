package pkg

import "fmt"

func updateModel(serverConfig *FinalServerConfig, model *FinalModel, inputData DataMap, whereArgs DataMap, requestedFields RequestedFieldMap) (interface{}, error) {
	dataLoaderWhereArgs := convertDataLoaderInput(model, whereArgs)
	fmt.Println("User Input:", inputData["user"])

	// Separate Linked vs Regular fields
	linkedInputFields := DataMap{}
	regularInputFields := DataMap{}
	for inputField, inputValue := range convertDataLoaderInput(model, inputData) {
		if value, isMap := inputValue.(DataMap); isMap {
			linkedInputFields[inputField] = value
		} else {
			regularInputFields[inputField] = inputValue
		}
	}

	// Update the regular fields
	dataLoaderNewItem, err := model.DataLoader.DataLoader.Update(model, regularInputFields, dataLoaderWhereArgs, requestedFields)
	if err != nil {
		return nil, err
	}

	// Update the linked fields
	for fieldName, fieldValue := range linkedInputFields {
		linkedField := model.Fields[fieldName].(*FinalLinkedField)
		nextModel := linkedField.LinkedModel
		nextInputData := fieldValue.(DataMap)
		convertedSourceField := model.Fields[linkedField.Field].(*FinalField).DataField
		notConvertedLinkedField := nextModel.Fields[linkedField.LinkedField].(*FinalField).Name
		nextWhereArgs := DataMap{
			notConvertedLinkedField: dataLoaderNewItem[convertedSourceField],
		}
		nextRequestedFields := requestedFields[fieldName].InnerFields.(RequestedFieldMap)

		_, err := updateModel(serverConfig, nextModel, nextInputData, nextWhereArgs, nextRequestedFields)
		if err != nil {
			return nil, err
		}
	}

	return getOneModel(serverConfig, model, whereArgs, requestedFields)
}
