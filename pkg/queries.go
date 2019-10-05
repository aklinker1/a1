package pkg

func getOneModel(serverConfig *FinalServerConfig, model *FinalModel, args DataMap, requestedFields RequestedFieldMap) (DataMap, error) {
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

	// Apply linked data
	err = applyLinks(serverConfig, model, requestedFields, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getMultipleModels(serverConfig *FinalServerConfig, model *FinalModel, args DataMap, requestedFields RequestedFieldMap) ([]DataMap, error) {
	// Get data
	dataLoaderItems, err := model.DataLoader.DataLoader.GetMultiple(
		model,
		convertDataLoaderInput(model, args),
		requestedFields,
	)
	if err != nil {
		return nil, err
	}

	// Apply linked data
	items := make([]DataMap, len(dataLoaderItems))
	for index, dataLoaderItem := range dataLoaderItems {
		items[index] = convertDataLoaderOutput(model, dataLoaderItem)
		err = applyLinks(serverConfig, model, requestedFields, dataLoaderItem)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func applyLinks(serverConfig *FinalServerConfig, model *FinalModel, requestedFields RequestedFieldMap, data DataMap) (err error) {
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
	// isAlreadyLinkedMap := map[string]*LinkedField{}
	// for fieldName, field := range model.Fields {
	// 	link := field.Linking

	// 	if link != nil {
	// 		// If there is a linking object directly on the model
	// 		isAlreadyLinkedMap[link.AccessedAs] = link
	// 		requestedType, areRequestingLinkedField := requestedFields[link.AccessedAs]
	// 		nextRequestedFields, isRequestedFieldMap := requestedType.(DataMap)
	// 		if areRequestingLinkedField && isRequestedFieldMap {
	// 			utils.Log("Linking %s.%s by %s.%s=%v", modelName, link.AccessedAs, link.ModelName, fieldName, data[fieldName])
	// 			linkedValue, err := selectOne(serverConfig, link.ModelName, serverConfig.Models[link.ModelName], data[fieldName], nextRequestedFields)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if len(linkedValue) != 0 {
	// 				data[link.AccessedAs] = linkedValue
	// 			}
	// 		}
	// 	}
	// }
	// for requestedField := range requestedFields {
	// 	nextRequestedFields, isLinkedObject := requestedFields[requestedField].(DataMap)
	// 	_, modelHasRequestedField := model.Fields[requestedField]
	// 	_, isAlreadyLinked := isAlreadyLinkedMap[requestedField]
	// 	if isLinkedObject && !modelHasRequestedField && !isAlreadyLinked {
	// 		for nextModelName, nextModel := range serverConfig.Models {
	// 			for fieldName, field := range nextModel.Fields {
	// 				if field.Linking != nil && field.Linking.ReverseAccessedAs == requestedField {
	// 					link := field.Linking
	// 					utils.Log("Reverse linking %s.%s by %s.%s=%v", link.ModelName, requestedField, nextModelName, fieldName, data[link.ForeignKey])
	// 					searchArgs := map[string]interface{}{}
	// 					items, _ := selectMultiple(serverConfig, nextModelName, nextModel, searchArgs, nextRequestedFields)
	// 					data[requestedField] = items
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return nil
}
