package pkg

func convertDataLoaderOutput(model *FinalModel, dataLoaderData DataMap) DataMap {
	mappedData := DataMap{}
	for dataLoaderField, value := range dataLoaderData {
		mappedField := model.DataFields[dataLoaderField]
		mappedData[mappedField] = value
	}
	return mappedData
}

func convertDataLoaderInput(model *FinalModel, args DataMap) DataMap {
	mappedData := DataMap{}
	for arg, value := range args {
		var mappedField string
		var finalValue interface{}
		switch field := model.Fields[arg].(type) {
		case *FinalField:
			mappedField = field.DataField
			finalValue = value
		case *FinalLinkedField:
			nextArgs := value.(DataMap)
			mappedField = field.Name
			finalValue = convertDataLoaderInput(field.LinkedModel, nextArgs)
		}
		mappedData[mappedField] = finalValue
	}
	return mappedData
}
