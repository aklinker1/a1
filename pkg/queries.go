package pkg

func getOneModel(serverConfig *FinalServerConfig, model *FinalModel, primaryKey interface{}, requestedFields RequestedFieldMap) (DataMap, error) {
	// Get data
	data, err := model.DataLoader.DataLoader.GetOne(model, primaryKey, requestedFields)
	if err != nil {
		return nil, err
	}

	// Apply linked data
	// err = applyLinks(data, serverConfig, modelName, model, requestedFields)
	// if err != nil {
	// 	return nil, err
	// }

	return data, nil
}

// func selectMultiple(serverConfig FinalServerConfig, model FinalModel, args DataMap, requestedFields RequestedFieldMap) ([]DataMap, error) {
// 	// Get data
// 	items, err := serverConfig.DatabaseDriver.SelectMultiple(model, args, requestedFields)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Apply linked data
// 	for _, item := range items {
// 		err = applyLinks(item, serverConfig, modelName, model, requestedFields)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return items, nil
// }

// func applyLinks(data DataMap, serverConfig FinalServerConfig, model FinalModel, requestedFields RequestedFieldMap) (err error) {
// 	isAlreadyLinkedMap := map[string]*LinkedField{}
// 	for fieldName, field := range model.Fields {
// 		link := field.Linking

// 		if link != nil {
// 			// If there is a linking object directly on the model
// 			isAlreadyLinkedMap[link.AccessedAs] = link
// 			requestedType, areRequestingLinkedField := requestedFields[link.AccessedAs]
// 			nextRequestedFields, isRequestedFieldMap := requestedType.(DataMap)
// 			if areRequestingLinkedField && isRequestedFieldMap {
// 				utils.Log("Linking %s.%s by %s.%s=%v", modelName, link.AccessedAs, link.ModelName, fieldName, data[fieldName])
// 				linkedValue, err := selectOne(serverConfig, link.ModelName, serverConfig.Models[link.ModelName], data[fieldName], nextRequestedFields)
// 				if err != nil {
// 					return err
// 				}
// 				if len(linkedValue) != 0 {
// 					data[link.AccessedAs] = linkedValue
// 				}
// 			}
// 		}
// 	}
// 	for requestedField := range requestedFields {
// 		nextRequestedFields, isLinkedObject := requestedFields[requestedField].(DataMap)
// 		_, modelHasRequestedField := model.Fields[requestedField]
// 		_, isAlreadyLinked := isAlreadyLinkedMap[requestedField]
// 		if isLinkedObject && !modelHasRequestedField && !isAlreadyLinked {
// 			for nextModelName, nextModel := range serverConfig.Models {
// 				for fieldName, field := range nextModel.Fields {
// 					if field.Linking != nil && field.Linking.ReverseAccessedAs == requestedField {
// 						link := field.Linking
// 						utils.Log("Reverse linking %s.%s by %s.%s=%v", link.ModelName, requestedField, nextModelName, fieldName, data[link.ForeignKey])
// 						searchArgs := map[string]interface{}{}
// 						items, _ := selectMultiple(serverConfig, nextModelName, nextModel, searchArgs, nextRequestedFields)
// 						data[requestedField] = items
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }
