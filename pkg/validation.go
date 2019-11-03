package pkg

import (
	"fmt"
	"os"
	"strings"
)

// Availability helpers for errors

func availableDataLoaders(serverConfig *FinalServerConfig) []string {
	availableOptions := []string{}
	for dataLoaderName := range serverConfig.DataLoaders {
		availableOptions = append(availableOptions, dataLoaderName)
	}
	return availableOptions
}

func availableModels(serverConfig *FinalServerConfig) []string {
	availableOptions := []string{}
	for modelName := range serverConfig.Models {
		availableOptions = append(availableOptions, modelName)
	}
	return availableOptions
}

func availableTypes(serverConfig *FinalServerConfig) []string {
	availableOptions := []string{}
	for typeName := range serverConfig.Types {
		availableOptions = append(availableOptions, typeName)
	}
	return availableOptions
}

func availableFields(model *FinalModel) []string {
	availableOptions := []string{}
	for fields := range model.Fields {
		availableOptions = append(availableOptions, fields)
	}
	return availableOptions
}

// Validators

func validateServerConfig(serverConfig *FinalServerConfig) []error {
	errors := []error{}
	// Basic Properties
	port := serverConfig.Port
	if port < 0 || port > 65535 {
		errors = append(errors, fmt.Errorf("Port (%d) must be between 0 and 65535", port))
	}
	endpoint := serverConfig.Endpoint
	if endpoint != "" && !strings.HasPrefix(endpoint, "/") {
		errors = append(errors, fmt.Errorf("Endpoint (%s) must start with a '/'", endpoint))
	}

	// GraphQL
	errors = append(errors, validateModels(serverConfig, serverConfig.Models)...)
	// TODO validate other stuff
	return errors
}

func validateModels(serverConfig *FinalServerConfig, models FinalModelMap) []error {
	errors := []error{}
	for _, model := range models {
		if model.PrimaryKey == "" && model.GraphQL.DisableGetOne == false {
			errors = append(errors, fmt.Errorf("Models[%s].PrimaryKey: Models must have a primary key, unless GraphQL.DisableGetOne is true", model.Name))
		}
		errors = append(errors, validateModelDataLoaderConfig(serverConfig, model, model.DataLoader)...)
		errors = append(errors, validateModelFields(serverConfig, model, model.Fields)...)
	}
	return errors
}

func validateModelDataLoaderConfig(serverConfig *FinalServerConfig, model *FinalModel, dataLoaderConfig FinalDataLoaderConfig) []error {
	errors := []error{}
	if dataLoaderConfig.DataLoader == nil {
		errors = append(errors, fmt.Errorf(
			"ServerConfig.Models[%s].DataLoader.Source: Cannot find data loader named '%s' (Available options: %v)",
			model.Name, dataLoaderConfig.Source, availableDataLoaders(serverConfig),
		))
	}
	if dataLoaderConfig.Group == "" {
		errors = append(errors, fmt.Errorf(
			"ServerConfig.Models[%s].DataLoader.Group: The data 'Group' property must be defined",
			model.Name,
		))
	}
	return errors
}

func validateModelFields(serverConfig *FinalServerConfig, model *FinalModel, fields FinalFieldMap) []error {
	errors := []error{}
	for fieldName, field := range fields {
		switch field.(type) {
		case *FinalField:
			regularField := field.(*FinalField)
			if regularField.Type == nil {
				errors = append(errors, fmt.Errorf(
					"ServerConfig.Models[%s].Fields[%s].Type: %s is not a valid type (Available options: %v)",
					model.Name, regularField.Name, regularField.TypeName, availableTypes(serverConfig),
				))
			}
		case *FinalVirtualField:
			virtualField := field.(*FinalVirtualField)
			if virtualField.Type == nil {
				errors = append(errors, fmt.Errorf(
					"ServerConfig.Models[%s].Fields[%s].Type: %s is not a valid type (Available options: %v)",
					model.Name, virtualField.Name, virtualField.TypeName, availableTypes(serverConfig),
				))
			}
		case *FinalLinkedField:
			linkedField := field.(*FinalLinkedField)
			if _, hasThisField := model.Fields[linkedField.Field]; !hasThisField {
				errors = append(errors, fmt.Errorf(
					"ServerConfig.Models[%s].Fields[%s].LinkedField.Field: %s is not a field on %s. (Available options: %v)",
					model.Name, linkedField.Name, linkedField.Field, model.Name, availableFields(model),
				))
				break
			}
			if linkedField.LinkedModel == nil {
				errors = append(errors, fmt.Errorf(
					"ServerConfig.Models[%s].Fields[%s].LinkedModel: %s is not a valid model name (Available options: %v)",
					model.Name, linkedField.Name, linkedField.LinkedModelName, availableModels(serverConfig),
				))
				break
			}
			if _, hasOtherField := linkedField.LinkedModel.Fields[linkedField.LinkedField]; !hasOtherField {
				errors = append(errors, fmt.Errorf(
					"ServerConfig.Models[%s].Fields[%s].LinkedField.LinkedField: %s is not a field on %s. (Available options: %v)",
					model.Name, linkedField.Name, linkedField.LinkedField, linkedField.LinkedModel.Name, availableFields(linkedField.LinkedModel),
				))
			}
		default:
			errors = append(errors, fmt.Errorf(
				"ServerConfig.Models[%s].Fields[%s]: %T is not an allowed type for a field",
				model.Name, fieldName, field,
			))
		}
	}
	return errors
}

func validateType(types FinalCustomTypeMap, typeName string) *FinalCustomType {
	tName := typeName
	if strings.HasSuffix(tName, "?") {
		tName = strings.ReplaceAll(tName, "?", "")
	}
	t, tExists := types[tName]
	if !tExists {
		fmt.Printf("Type '%s' is not a type\n", typeName)
		os.Exit(1)
	}
	return t
}
