package pkg

func updateModel(serverConfig *FinalServerConfig, model *FinalModel, dataChanges DataMap, whereArgs DataMap, requestedFields RequestedFieldMap) (interface{}, error) {
	// Update data
	dataLoaderData, err := model.DataLoader.DataLoader.Update(
		model,
		convertDataLoaderInput(model, dataChanges),
		convertDataLoaderInput(model, whereArgs),
		requestedFields,
	)
	updatedData := convertDataLoaderOutput(model, dataLoaderData)
	if err != nil {
		return nil, err
	}

	// Apply linked data & compute virtual fields
	err = applyLinks(serverConfig, model, requestedFields, updatedData)
	if err != nil {
		return nil, err
	}
	err = computeVirtualFields(serverConfig, model, requestedFields, updatedData)
	if err != nil {
		return nil, err
	}

	return updatedData, nil
}
