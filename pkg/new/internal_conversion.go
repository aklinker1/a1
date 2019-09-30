package new

import (
	"fmt"
	"os"

	graphql "github.com/graphql-go/graphql"
)

func convertTypes(types CustomTypeMap) FinalCustomTypeMap {
	finalTypes := FinalCustomTypeMap{
		"Bool": &FinalCustomType{
			Name:        "Bool",
			GraphQLType: graphql.Boolean,
		},
		"Int": &FinalCustomType{
			Name:        "Int",
			GraphQLType: graphql.Int,
		},
		// "Long": &FinalCustomType{
		// 	Name:        "Long",
		// 	GraphQLType: graphql.Long,
		// },
		"Float": &FinalCustomType{
			Name:        "Float",
			GraphQLType: graphql.Float,
		},
		// "Double": &FinalCustomType{
		// 	Name:        "Double",
		// 	GraphQLType: graphql.Double,
		// },
		"String": &FinalCustomType{
			Name:        "String",
			GraphQLType: graphql.String,
		},
		"ID": &FinalCustomType{
			Name:        "ID",
			GraphQLType: graphql.ID,
		},
		"Date": &FinalCustomType{
			Name:        "Date",
			GraphQLType: graphql.DateTime,
		},
	}
	for customTypeName, customType := range types {
		finalTypes[customTypeName] = &FinalCustomType{
			Name:        customTypeName,
			Description: customType.Description,
			ToJSON:      customType.ToJSON,
			FromJSON:    customType.FromJSON,
			FromLiteral: customType.FromLiteral,
		}
	}
	return finalTypes
}

func convertDataLoader(name string, dataLoader DataLoader) FinalDataLoader {
	return FinalDataLoader{
		Name:        name,
		Connect:     dataLoader.Connect,
		GetOne:      dataLoader.GetOne,
		GetMultiple: dataLoader.GetMultiple,
	}
}

func convertModelsToLinkedModels(types FinalCustomTypeMap, models ModelMap) ModelMap {
	linkedModels := ModelMap{}
	for modelName, model := range models {
		for _, field := range model.Fields {
			linkedField, isLinkedField := field.(LinkedField)
			if isLinkedField {
				linkedModelName := linkedField.getCustomModelName(modelName)
				linkedModel := models[linkedField.LinkedModel]
				fmt.Println()
				fmt.Println("Model:", linkedModelName)
				linkedModels[linkedModelName] = Model{
					Description: linkedModel.Description,
					Extends:     linkedModel.Extends,
					Fields:      removeFieldsForLink(linkedModel.Fields, linkedField.Field),
					DataLoader:  linkedModel.DataLoader,
					GraphQL:     linkedModel.GraphQL,
				}
			}
		}
	}
	return linkedModels
}

func removeFieldsForLink(fields FieldMap, fieldToRemove string) FieldMap {
	newFields := FieldMap{}
	for fieldName, field := range fields {
		newFields[fieldName] = field
	}
	for linkedFieldName, field := range newFields {
		linkedField, isLinkedField := field.(LinkedField)
		if isLinkedField && linkedField.LinkedField == fieldToRemove {
			delete(newFields, linkedFieldName)
			fmt.Println("Removing Linked Field from child: ", linkedFieldName)
		}
	}
	return newFields
}

type linkedFieldAndModel struct {
	LinkedField  *FinalLinkedField
	Model        *FinalModel
	KeyToExclude string
}

func getLinkedFieldsToAppend(models ModelMap, allModels FinalModelMap) []linkedFieldAndModel {
	// Generate a list of models -> field to add
	fieldsToLink := []linkedFieldAndModel{}
	for modelName, model := range models {
		finalModel := allModels[modelName]
		for fieldName, field := range model.Fields {
			linkedField, isLinkedField := field.(LinkedField)
			if isLinkedField {
				newLinkedField := linkedFieldAndModel{
					LinkedField:  convertLinkedField(allModels, fieldName, linkedField, finalModel),
					Model:        finalModel,
					KeyToExclude: linkedField.Field,
				}
				fieldsToLink = append(fieldsToLink, newLinkedField)
			}
		}
	}
	return fieldsToLink
}

func convertModelToFinalWithoutLinkedFields(dataLoaders FinalDataLoaderMap, types FinalCustomTypeMap, modelName string, model Model) *FinalModel {
	fields := FinalFieldMap{}
	var primaryKey string
	primaryKeyCount := 0
	for fieldName, field := range model.Fields {
		if _, isLinkedField := field.(LinkedField); !isLinkedField {
			fields[fieldName] = convertNonLinkedField(types, fieldName, field)
		}
		if regularField, isRegularField := field.(Field); isRegularField && regularField.PrimaryKey {
			primaryKeyCount++
			primaryKey = fieldName
		}
	}
	if primaryKeyCount > 1 {
		fmt.Printf("Cannot have more than 1 primary key defined on a model. You have %d on '%s' \n", primaryKeyCount, modelName)
		os.Exit(1)
	}
	return &FinalModel{
		Name:        modelName,
		Extended:    model.Extends,
		Description: model.Description,
		DataLoader:  convertDataLoaderConfig(dataLoaders, model.DataLoader),
		Fields:      fields,
		GraphQL:     convertGraphQLConfig(model.GraphQL),
		PrimaryKey:  primaryKey,
	}
}

func extendModels(models FinalModelMap) {
	for _, model := range models {
		if model.Extended != "" {
			parentModelPointer, parentModelExists := models[model.Extended]
			if parentModelExists {
				models[model.Name] = model.extends(parentModelPointer)
			}
		}
	}
}

func convertLinkedField(models FinalModelMap, fieldName string, field LinkedField, parentModel *FinalModel) *FinalLinkedField {
	linkedModelName := field.getCustomModelName(parentModel.Name)
	linkedModel := models[linkedModelName]

	return &FinalLinkedField{
		Name:              fieldName,
		Description:       field.Description,
		DeprecationReason: field.DeprecationReason,
		LinkedModelName:   linkedModelName,
		LinkedModel:       linkedModel,
		Field:             field.Field,
		LinkedField:       field.LinkedField,
	}
}

func convertDataLoaderConfig(dataLoaders FinalDataLoaderMap, config DataLoaderConfig) FinalDataLoaderConfig {
	return FinalDataLoaderConfig{
		DataLoader: dataLoaders[config.Source],
		Source:     config.Source,
		Group:      config.Group,
	}
}

func convertGraphQLConfig(config GraphQLConfig) FinalGraphQLConfig {
	return FinalGraphQLConfig{
		DisableSelectOne:      config.DisableSelectOne,
		DisableSelectMultiple: config.DisableSelectMultiple,
		DisableCreate:         config.DisableCreate,
		DisableUpdate:         config.DisableUpdate,
		DisableDelete:         config.DisableDelete,
	}
}

func convertNonLinkedField(types FinalCustomTypeMap, fieldName string, field GraphQLField) FinalGraphQLField {
	switch field.(type) {
	case string:
		typeName := field.(string)
		return &FinalField{
			Name:              fieldName,
			Description:       "",
			DeprecationReason: "",
			TypeName:          typeName,
			Type:              types[typeName],
			Hidden:            false,
			PrimaryKey:        false,
			DataField:         fieldName,
		}
	case Field:
		regularField := field.(Field)
		dataFieldName := fieldName
		if regularField.DataField != "" {
			dataFieldName = regularField.DataField
		}
		return &FinalField{
			Name:              fieldName,
			Description:       regularField.Description,
			DeprecationReason: regularField.DeprecationReason,
			TypeName:          regularField.Type,
			Type:              types[regularField.Type],
			Hidden:            regularField.Hidden,
			PrimaryKey:        regularField.PrimaryKey,
			DataField:         dataFieldName,
		}
	case VirtualField:
		virtualField := field.(VirtualField)
		return &FinalVirtualField{
			Name:              fieldName,
			Description:       virtualField.Description,
			DeprecationReason: virtualField.DeprecationReason,
			TypeName:          virtualField.Type,
			Type:              types[virtualField.Type],
		}
	}
	return nil
}

func convertEnum(enumName string, enum Enum) *FinalCustomType {
	return &FinalCustomType{
		GraphQLType: graphql.NewEnum(graphql.EnumConfig{
			Name:        enumName,
			Description: enum.Description,
			Values:      convertEnumValues(enum.Values),
		}),
		Name:        enumName,
		Description: enum.Description,
	}
}

func convertEnumValues(values EnumValueMap) graphql.EnumValueConfigMap {
	graphqlValues := graphql.EnumValueConfigMap{}
	for valueName, value := range values {
		graphqlValues[valueName] = &graphql.EnumValueConfig{
			Value:             value.Value,
			DeprecationReason: value.DeprecationReason,
			Description:       value.Description,
		}
	}
	return graphqlValues
}
