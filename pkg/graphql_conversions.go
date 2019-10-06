package pkg

import (
	"fmt"

	utils "github.com/aklinker1/a1/pkg/utils"
	errors "github.com/go-errors/errors"
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
		name := "Input_" + modelName
		inputs[name] = graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        name,
			Description: model.Description,
			Fields:      inputFieldsWithoutLinked(types, model.Fields),
		})
	}

	return inputs
}

func inputFieldsWithoutLinked(types graphql.TypeMap, fields FinalFieldMap) graphql.InputObjectConfigFieldMap {
	inputFields := graphql.InputObjectConfigFieldMap{}
	for fieldName, field := range fields {
		if regularField, isRegularField := field.(*FinalField); isRegularField {
			if !regularField.Hidden {
				inputFields[fieldName] = &graphql.InputObjectFieldConfig{
					Type:        types[regularField.Type.Name],
					Description: regularField.Description,
				}
			}
		}
	}
	return inputFields
}

func addLinksToInputObjects(inputs map[string]*graphql.InputObject, models FinalModelMap) {
	for _, model := range models {
		for _, field := range model.Fields {
			linkedField, isLinkedField := field.(*FinalLinkedField)
			if isLinkedField {
				linkedInputModelName := "Input_" + linkedField.LinkedModelName
				inputs["Input_"+model.Name].AddFieldConfig(linkedField.Name, &graphql.InputObjectFieldConfig{
					Type:        inputs[linkedInputModelName],
					Description: linkedField.Description,
				})
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
		case *FinalField:
			regularField := field.(*FinalField)
			if !regularField.Hidden {
				outputFields[fieldName] = &graphql.Field{
					Name:              fieldName,
					Type:              types[regularField.Type.Name],
					Description:       regularField.Description,
					DeprecationReason: regularField.DeprecationReason,
				}
			}
		case *FinalVirtualField:
			virtualField := field.(*FinalVirtualField)
			outputFields[fieldName] = &graphql.Field{
				Name:              fieldName,
				Type:              types[virtualField.Type.Name],
				Description:       virtualField.Description,
				DeprecationReason: virtualField.DeprecationReason,
			}
		}
	}
	return outputFields
}

func addLinksToOutputObjects(outputs map[string]*graphql.Object, models FinalModelMap) {
	for _, model := range models {
		for _, field := range model.Fields {
			linkedField, isLinkedField := field.(*FinalLinkedField)
			if isLinkedField {
				if linkedField.Type == OneToMany {
					outputs[model.Name].AddFieldConfig(linkedField.Name, &graphql.Field{
						Name:              linkedField.Name,
						DeprecationReason: linkedField.DeprecationReason,
						Type:              graphql.NewList(outputs[linkedField.LinkedModelName]),
						Description:       linkedField.Description,
					})
				} else {
					outputs[model.Name].AddFieldConfig(linkedField.Name, &graphql.Field{
						Name:              linkedField.Name,
						DeprecationReason: linkedField.DeprecationReason,
						Type:              outputs[linkedField.LinkedModelName],
						Description:       linkedField.Description,
					})
				}
			}
		}
	}
}

// RESOLVERS

func (resolver Resolvable) graphqlResolver() func(params graphql.ResolveParams) (interface{}, error) {
	isVerbose := utils.IsVerbose()
	return func(params graphql.ResolveParams) (interface{}, error) {
		utils.Log("\n  %s(args: %v)", resolver.Name, resolver.Arguments)
		// Check Authorization
		// myUser, err := middleware.Authorize(params.Context, function.AuthRequired)
		// if err != nil {
		// 	return 401, err
		// }

		// // Check Role Level
		// if myUser != nil {
		// 	err = middleware.CheckRole(myUser, function.MinRole)
		// 	if err != nil {
		// 		return 403, err
		// 	}
		// }

		// Get argument map
		args := params.Args

		// Get field map
		fields, err := GetRequestedFields(serverConfig, resolver.Model, params)
		if isVerbose {
			utils.Log("Field Map:")
			printRequestedFields(fields, 2)
		}

		// Call Resolver
		result, err := resolver.Resolver(args, fields)
		if err != nil {
			stack := errors.New(err).ErrorStack()
			fmt.Println("STACK", stack)
			return nil, err
		}
		resultMap, isMap := result.(DataMap)
		if isMap {
			utils.Log("Resolved JSON Object: %v", resultMap)
			if len(resultMap) == 0 {
				return nil, nil
			}
			return result, nil
		}
		resultArray, isArray := result.([]DataMap)
		if isArray {
			utils.Log("Resolved JSON Array[%d]: %v", len(resultArray), resultArray)
			if len(resultArray) == 0 {
				return []interface{}{}, nil
			}
			return resultArray, nil
		}
		return result, nil
	}
}

func graphqlArguments(arguments []Argument, allTypes graphql.TypeMap) graphql.FieldConfigArgument {
	config := graphql.FieldConfigArgument{}
	for _, argument := range arguments {
		config[argument.Name] = &graphql.ArgumentConfig{
			Type:         allTypes[argument.Type],
			DefaultValue: argument.DefaultValue,
			Description:  argument.Description,
		}
	}
	return config
}

func (resolver Resolvable) graphqlResolverEntry(allTypes graphql.TypeMap) *graphql.Field {
	var returnType graphql.Output = allTypes[resolver.Model.Name]
	if resolver.ResturnsList {
		returnType = graphql.NewList(allTypes[resolver.Model.Name])
	}
	return &graphql.Field{
		Name:        resolver.Name,
		Description: resolver.Description,
		Args:        graphqlArguments(resolver.Arguments, allTypes),
		Type:        returnType,
		Resolve:     resolver.graphqlResolver(),
	}
}
