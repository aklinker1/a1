package pkg

import (
	utils "github.com/aklinker1/a1/pkg/utils"
)

func applyLinks(data StringMap, serverConfig ServerConfig, modelName string, model Model, requestedFields StringMap) (err error) {
	isAlreadyLinkedMap := map[string]*LinkedField{}
	for fieldName, field := range model.Fields {
		link := field.Linking

		if link != nil {
			// If there is a linking object directly on the model
			isAlreadyLinkedMap[link.AccessedAs] = link
			requestedType, areRequestingLinkedField := requestedFields[link.AccessedAs]
			nextRequestedFields, isRequestedFieldMap := requestedType.(StringMap)
			if areRequestingLinkedField && isRequestedFieldMap {
				utils.Log("Linking %s.%s by %s.%s=%v", modelName, link.AccessedAs, link.ModelName, fieldName, data[fieldName])
				linkedValue, err := selectOne(serverConfig, link.ModelName, serverConfig.Models[link.ModelName], data[fieldName], nextRequestedFields)
				if err != nil {
					return err
				}
				if len(linkedValue) != 0 {
					data[link.AccessedAs] = linkedValue
				}
			}
		}
	}
	for requestedField := range requestedFields {
		nextRequestedFields, isLinkedObject := requestedFields[requestedField].(StringMap)
		_, modelHasRequestedField := model.Fields[requestedField]
		_, isAlreadyLinked := isAlreadyLinkedMap[requestedField]
		if isLinkedObject && !modelHasRequestedField && !isAlreadyLinked {
			for nextModelName, nextModel := range serverConfig.Models {
				for fieldName, field := range nextModel.Fields {
					if field.Linking != nil && field.Linking.ReverseAccessedAs == requestedField {
						link := field.Linking
						utils.Log("Reverse linking %s.%s by %s.%s=%v", link.ModelName, requestedField, nextModelName, fieldName, data[link.ForeignKey])
						searchArgs := map[string]interface{}{}
						items, _ := selectMultiple(serverConfig, nextModelName, nextModel, searchArgs, nextRequestedFields)
						data[requestedField] = items
					}
				}
			}
		}
	}
	return nil
}

func selectOne(serverConfig ServerConfig, modelName string, model Model, primaryKey interface{}, requestedFields StringMap) (StringMap, error) {
	// Get data
	data, err := serverConfig.DatabaseDriver.SelectOne(model, primaryKey, requestedFields)
	if err != nil {
		return nil, err
	}

	// Apply linked data
	err = applyLinks(data, serverConfig, modelName, model, requestedFields)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func selectOneQuery(modelName string, model Model, serverConfig ServerConfig) *Resolvable {
	return &Resolvable{
		Model:     model,
		ModelName: modelName,
		Name:      utils.LowerFirstChar(modelName),
		Arguments: []Argument{
			Argument{
				Name: model.PrimaryKey,
				Type: "Int",
			},
		},
		Resolver: func(args StringMap, requestedFields StringMap) (interface{}, error) {
			return selectOne(serverConfig, modelName, model, args[model.PrimaryKey], requestedFields)
		},
	}
}

func selectMultiple(serverConfig ServerConfig, modelName string, model Model, args StringMap, requestedFields StringMap) ([]StringMap, error) {
	// Get data
	items, err := serverConfig.DatabaseDriver.SelectMultiple(model, args, requestedFields)
	if err != nil {
		return nil, err
	}

	// Apply linked data
	for _, item := range items {
		err = applyLinks(item, serverConfig, modelName, model, requestedFields)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func selectMultipleQuery(modelName string, model Model, serverConfig ServerConfig) *Resolvable {
	return &Resolvable{
		Model:        model,
		ModelName:    modelName,
		Name:         utils.AddS(utils.LowerFirstChar(modelName)),
		ResturnsList: true,
		Resolver: func(args StringMap, requestedFields StringMap) (interface{}, error) {
			return selectMultiple(serverConfig, modelName, model, args, requestedFields)
		},
	}
}
