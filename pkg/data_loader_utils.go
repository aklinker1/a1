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
		mappedField := model.Fields[arg].(*FinalField).DataField
		mappedData[mappedField] = value
	}
	return mappedData
}
