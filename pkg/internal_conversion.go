package pkg

import (
	"fmt"
	"os"

	utils "github.com/aklinker1/a1/pkg/utils"
	graphql "github.com/graphql-go/graphql"
)

func convertTypes(customTypes CustomTypeMap) FinalCustomTypeMap {
	finalTypes := FinalCustomTypeMap{}

	utils.LogWhite("[Builtin Types - %d]", len(builtinTypes))
	for builtinTypeName, builtinType := range builtinTypes {
		utils.Log("  - %s", builtinTypeName)
		finalTypes[builtinTypeName] = convertType(builtinTypeName, builtinType)
	}
	utils.LogWhite("[Custom Types - %d]", len(customTypes))
	for customTypeName, customType := range customTypes {
		utils.Log("  - %s", customTypeName)
		finalTypes[customTypeName] = convertType(customTypeName, customType)
	}

	return finalTypes
}

func convertType(name string, customType CustomType) *FinalCustomType {
	if customType.GraphQLType != nil {
		return &FinalCustomType{
			Name:        name,
			GraphQLType: customType.GraphQLType,
		}
	}
	return &FinalCustomType{
		Name:        name,
		Description: customType.Description,
		ToJSON:      customType.ToJSON,
		FromJSON:    customType.FromJSON,
		FromLiteral: customType.FromLiteral,
	}
}

func convertDataLoader(name string, dataLoader DataLoader) FinalDataLoader {
	return FinalDataLoader{
		Name:        name,
		Connect:     dataLoader.Connect,
		GetOne:      dataLoader.GetOne,
		GetMultiple: dataLoader.GetMultiple,
		Update:      dataLoader.Update,
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
	dataFields := StringMap{}
	var primaryKey string
	primaryKeyCount := 0
	for fieldName, field := range model.Fields {
		if _, isLinkedField := field.(LinkedField); !isLinkedField {
			nonLinkedField := convertNonLinkedField(types, fieldName, field)
			fields[fieldName] = nonLinkedField
			if regularField, isRegularField := nonLinkedField.(*FinalField); isRegularField {
				dataFields[regularField.DataField] = fieldName
			}
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
		DataFields:  dataFields,
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
		IsNullable:        field.IsNullable,
		LinkedModelName:   linkedModelName,
		LinkedModel:       linkedModel,
		Type:              field.Type,
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
		DisableGetOne:      config.DisableSelectOne,
		DisableGetMultiple: config.DisableSelectMultiple,
		DisableCreate:      config.DisableCreate,
		DisableUpdate:      config.DisableUpdate,
		DisableDelete:      config.DisableDelete,
	}
}

func convertNonLinkedField(types FinalCustomTypeMap, fieldName string, field GraphQLField) FinalGraphQLField {
	switch field.(type) {
	case string:
		typeName := field.(string)
		finalType := validateType(types, typeName)
		return &FinalField{
			Name:              fieldName,
			Description:       "",
			DeprecationReason: "",
			TypeName:          typeName,
			Type:              finalType,
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
		finalType := validateType(types, regularField.Type)
		return &FinalField{
			Name:              fieldName,
			Description:       regularField.Description,
			DeprecationReason: regularField.DeprecationReason,
			TypeName:          regularField.Type,
			Type:              finalType,
			Hidden:            regularField.Hidden,
			PrimaryKey:        regularField.PrimaryKey,
			DataField:         dataFieldName,
		}
	case VirtualField:
		virtualField := field.(VirtualField)
		finalType := validateType(types, virtualField.Type)
		return &FinalVirtualField{
			Name:              fieldName,
			Description:       virtualField.Description,
			DeprecationReason: virtualField.DeprecationReason,
			TypeName:          virtualField.Type,
			Type:              finalType,
			RequiredFields:    virtualField.RequiredFields,
			Compute:           virtualField.Compute,
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

func generateQueriesForModel(serverConfig *FinalServerConfig, model *FinalModel) []Resolvable {
	results := []Resolvable{}
	if !model.GraphQL.DisableGetOne {
		results = append(results, GetOneQuery(serverConfig, model))
	}
	if !model.GraphQL.DisableGetMultiple {
		results = append(results, GetMultipleQuery(serverConfig, model))
	}
	return results
}

func generateMutationsForModels(serverConfig *FinalServerConfig, model *FinalModel) []Resolvable {
	results := []Resolvable{}
	if !model.GraphQL.DisableCreate {
		// results = append(results, GetOneQuery(serverConfig, model))
	}
	if !model.GraphQL.DisableUpdate {
		results = append(results, UpdateMutation(serverConfig, model))
	}
	if !model.GraphQL.DisableDelete {
		// results = append(results, GetMultipleQuery(serverConfig, model))
	}
	return results
}
