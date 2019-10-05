package pkg

func getOneModel(
	serverConfig *FinalServerConfig,
	model *FinalModel,
	args DataMap,
	requestedFields RequestedFieldMap,
) (DataMap, error) {
	// Get data
	dataLoaderData, err := model.DataLoader.DataLoader.GetOne(
		model,
		convertDataLoaderInput(model, args),
		requestedFields,
	)
	data := convertDataLoaderOutput(model, dataLoaderData)
	if err != nil {
		return nil, err
	}

	// Apply linked data & compute virtual fields
	err = applyLinks(serverConfig, model, requestedFields, data)
	if err != nil {
		return nil, err
	}
	err = computeVirtualFields(serverConfig, model, requestedFields, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getMultipleModels(
	serverConfig *FinalServerConfig,
	model *FinalModel,
	args DataMap,
	requestedFields RequestedFieldMap,
) ([]DataMap, error) {
	// Get data
	dataLoaderItems, err := model.DataLoader.DataLoader.GetMultiple(
		model,
		convertDataLoaderInput(model, args),
		requestedFields,
	)
	if err != nil {
		return nil, err
	}

	items := make([]DataMap, len(dataLoaderItems))
	for index, dataLoaderItem := range dataLoaderItems {
		item := convertDataLoaderOutput(model, dataLoaderItem)

		// Apply linked data & compute virtual fields
		err = applyLinks(serverConfig, model, requestedFields, item)
		if err != nil {
			return nil, err
		}
		err = computeVirtualFields(serverConfig, model, requestedFields, item)
		if err != nil {
			return nil, err
		}

		items[index] = item
	}

	return items, nil
}

func applyLinks(
	serverConfig *FinalServerConfig,
	model *FinalModel,
	requestedFields RequestedFieldMap,
	data DataMap,
) (err error) {
	for _, requestedField := range requestedFields {
		if linkedField, ok := requestedField.Field.(*FinalLinkedField); ok {
			nextRequestedFields := requestedFields[linkedField.Name].InnerFields.(RequestedFieldMap)
			switch linkedField.Type {
			case OneToOne:
				fallthrough
			case ManyToOne:
				args := DataMap{
					linkedField.LinkedField: data[linkedField.Field],
				}
				innerData, err := getOneModel(
					serverConfig,
					linkedField.LinkedModel,
					args,
					nextRequestedFields,
				)
				if err != nil {
					return err
				}
				data[linkedField.Name] = innerData
			case OneToMany:
				args := DataMap{
					linkedField.LinkedField: data[linkedField.Field],
				}
				innerData, err := getMultipleModels(
					serverConfig,
					linkedField.LinkedModel,
					args,
					nextRequestedFields,
				)
				if err != nil {
					return err
				}
				data[linkedField.Name] = innerData
			}
		}
	}
	return nil
}

func computeVirtualFields(
	serverConfig *FinalServerConfig,
	model *FinalModel,
	requestedFields RequestedFieldMap,
	data DataMap,
) (err error) {
	for requestedFieldName, requestedField := range requestedFields {
		if virtualField, ok := requestedField.Field.(*FinalVirtualField); ok {
			data[requestedFieldName], err = virtualField.Compute(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
