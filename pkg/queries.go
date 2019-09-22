package pkg

import (
	utils "github.com/aklinker1/a1/pkg/utils"
)

func applyLinks(data StringMap, serverConfig ServerConfig, modelName string, model Model, requestedFields StringMap) (err error) {
	for fieldName, field := range model.Fields {
		link := field.Linking

		if link != nil {
			// If there is a linking object directly on the model
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
		_, modelHasRequestedField := model.Fields[requestedField]
		if !modelHasRequestedField {
			utils.Log("Field not mapped (%v): %v", model.Fields, requestedField)
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
		Returns:   model,
		Arguments: []Argument{
			Argument{
				Name: model.PrimaryKey,
				Type: "Int",
			},
		},
		Resolver: func(args StringMap, fields StringMap) (StringMap, error) {
			return selectOne(serverConfig, modelName, model, args[model.PrimaryKey], fields)
		},
	}
}
