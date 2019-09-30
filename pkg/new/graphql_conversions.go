package new

import (
	graphql "github.com/graphql-go/graphql"
)

// TYPES

func getGraphqlTypes(types FinalCustomTypeMap) graphql.TypeMap {
	output := graphql.TypeMap{}

	for _, customType := range types {
		if customType.GraphQLType != nil {
			output[customType.Name] = customType.GraphQLType
		} else {
			output[customType.Name] = graphql.NewScalar(graphql.ScalarConfig{
				Name:         customType.Name,
				Description:  customType.Description,
				Serialize:    customType.ToJSON,
				ParseValue:   customType.FromJSON,
				ParseLiteral: customType.FromLiteral,
			})
		}
	}
	return output
}

// INPUTS

func inputModelsWithoutLinkedFields(types graphql.TypeMap, models FinalModelMap) map[string]*graphql.InputObject {
	inputs := map[string]*graphql.InputObject{}

	for modelName, model := range models {
		inputs[modelName] = graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        modelName,
			Description: model.Description,
			Fields:      inputFieldsWithoutLinked(types, model.Fields),
		})
	}

	return inputs
}

func inputFieldsWithoutLinked(types graphql.TypeMap, fields FinalFieldMap) graphql.Fields {
	inputFields := graphql.Fields{}
	for fieldName, field := range fields {
		switch field.(type) {
		case Field:
			regularField := field.(Field)
			inputFields[fieldName] = &graphql.Field{
				Name:              fieldName,
				Type:              types[regularField.Type],
				Description:       regularField.Description,
				DeprecationReason: regularField.DeprecationReason,
			}
		case VirtualField:
			virtualField := field.(VirtualField)
			inputFields[fieldName] = &graphql.Field{
				Name:              fieldName,
				Type:              types[virtualField.Type],
				Description:       virtualField.Description,
				DeprecationReason: virtualField.DeprecationReason,
			}
		}
	}
	return inputFields
}

func addLinksToInputObjects(inputs map[string]*graphql.InputObject, models FinalModelMap) {
	var count int
	for _, model := range models {
		for _, field := range model.Fields {
			linkedField, isLinkedField := field.(FinalLinkedField)
			if isLinkedField {
				inputs[model.Name].AddFieldConfig(linkedField.Name, &graphql.InputObjectFieldConfig{
					Type:        inputs[model.Name],
					Description: linkedField.Description,
				})
				count++
			}
		}
	}
}

// OUTPUTS

func outputModelsWithoutLinkedFields(types graphql.TypeMap, models FinalModelMap) map[string]*graphql.Object {
	outputs := map[string]*graphql.Object{}

	for modelName, model := range models {
		outputs[modelName] = graphql.NewObject(graphql.ObjectConfig{
			Name:        modelName,
			Description: model.Description,
			Fields:      outputFieldsWithoutLinked(types, model.Fields),
		})
	}

	return outputs
}

func outputFieldsWithoutLinked(types graphql.TypeMap, fields FinalFieldMap) graphql.Fields {
	outputFields := graphql.Fields{}
	for fieldName, field := range fields {
		switch field.(type) {
		case Field:
			regularField := field.(Field)
			outputFields[fieldName] = &graphql.Field{
				Name:              fieldName,
				Type:              types[regularField.Type],
				Description:       regularField.Description,
				DeprecationReason: regularField.DeprecationReason,
			}
		case VirtualField:
			virtualField := field.(VirtualField)
			outputFields[fieldName] = &graphql.Field{
				Name:              fieldName,
				Type:              types[virtualField.Type],
				Description:       virtualField.Description,
				DeprecationReason: virtualField.DeprecationReason,
			}
		}
	}
	return outputFields
}

func addLinksToOutputObjects(outputs map[string]*graphql.Object, models FinalModelMap) {
	var count int
	for _, model := range models {
		for _, field := range model.Fields {
			linkedField, isLinkedField := field.(FinalLinkedField)
			if isLinkedField {
				outputs[model.Name].AddFieldConfig(linkedField.Name, &graphql.Field{
					Name:              linkedField.Name,
					DeprecationReason: linkedField.DeprecationReason,
					Type:              outputs[model.Name],
					Description:       linkedField.Description,
				})
				count++
			}
		}
	}
}
