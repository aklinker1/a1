package pkg

func applyLinks(data StringMap, serverConfig ServerConfig, model Model, requestedFields StringMap) (err error) {
	for fieldName, field := range model.Fields {
		link := field.Linking
		if link != nil {
			requestedType, areRequestingLinkedField := requestedFields[link.AccessedAs]
			nextRequestedFields, isRequestedFieldMap := requestedType.(StringMap)
			if areRequestingLinkedField && isRequestedFieldMap {
				var linkedValue StringMap
				switch link.Type {
				case OneToOne:
					Log("Linking %s by %s=%v", link.ModelName, link.AccessedAs, data[fieldName])
					linkedValue, err = selectOne(serverConfig, serverConfig.Models[link.ModelName], data[fieldName], nextRequestedFields)
					if err != nil {
						return err
					}
				case OneToMany:
					// TODO
					linkedValue = StringMap{}
				}
				if len(linkedValue) != 0 {
					data[link.AccessedAs] = linkedValue
				}
			}
		}
	}
	return nil
}

func selectOne(serverConfig ServerConfig, model Model, primaryKey interface{}, requestedFields StringMap) (StringMap, error) {
	// TODO: Mocking for now, remove
	requestedFields = StringMap{
		"id":      "ID",
		"title":   "String",
		"_userId": "ID",
		"user": StringMap{
			"id":       "ID",
			"username": "String",
		},
	}

	// Get data
	data, err := serverConfig.DatabaseDriver.SelectOne(model, primaryKey, requestedFields)
	if err != nil {
		return nil, err
	}

	// Apply linked data
	err = applyLinks(data, serverConfig, model, requestedFields)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func selectOneQuery(modelName string, model Model, serverConfig ServerConfig) *Resolvable {
	return &Resolvable{
		Model:     model,
		ModelName: modelName,
		Name:      lowerFirstChar(modelName),
		Returns:   model,
		Arguments: []Argument{
			Argument{
				Name: model.PrimaryKey,
				Type: "Int",
			},
		},
		Resolver: func(args StringMap, fields StringMap) (StringMap, error) {
			return selectOne(serverConfig, model, args[model.PrimaryKey], fields)
		},
	}
}
