package new

import (
	"fmt"
	"os"
)

// Start -
func Start(serverConfig ServerConfig) {
	// Load ENV Variables

	// Parse Server Config
	finalServerconfig, errors := parseServerConfig(serverConfig)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	// Validate FinalServerConfig
	errors = validateServerConfig(finalServerconfig)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	// Connect the data loaders
	for _, dataLoader := range finalServerconfig.DataLoaders {
		err := dataLoader.Connect()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Setup the GraphQL schema
}

func parseServerConfig(serverConfig ServerConfig) (FinalServerConfig, []error) {
	// Get final data loaders
	dataLoaders := FinalDataLoaderMap{}
	for dataLoaderName, dataLoader := range serverConfig.DataLoaders {
		finalDataLoader := convertDataLoader(dataLoaderName, dataLoader)
		dataLoaders[dataLoaderName] = &finalDataLoader
	}

	// Get final types
	types := convertTypes(serverConfig.Types)

	// Append enums to types
	for enumName, enum := range serverConfig.Enums {
		types[enumName] = convertEnum(enumName, enum)
	}

	// Generate Linked Model types
	fmt.Printf("%d base models\n", len(serverConfig.Models))
	for modelName := range serverConfig.Models {
		fmt.Printf("- %s\n", modelName)
	}

	linkedModels := convertModelsToLinkedModels(types, serverConfig.Models)
	fmt.Printf("\nGenerated %d linked models\n", len(linkedModels))
	for linkedModelName, linkedModel := range linkedModels {
		fmt.Printf("- %s\n", linkedModelName)
		serverConfig.Models[linkedModelName] = linkedModel
	}
	fmt.Printf("\n%d models total\n", len(serverConfig.Models))

	// Get Models (without linked fields)
	models := FinalModelMap{}
	for modelName, model := range serverConfig.Models {
		models[modelName] = convertModelToFinalWithoutLinkedFields(dataLoaders, types, modelName, model)
	}

	// Copy model properties to extended models
	extendModels(models)

	// Append linked fields to models
	fieldsToAppend := getLinkedFieldsToAppend(serverConfig.Models, models)
	for _, modelAndField := range fieldsToAppend {
		modelAndField.Model.Fields[modelAndField.LinkedField.Name] = modelAndField.LinkedField
	}
	fmt.Printf("\nAdded %d linked fields to the models\n", len(fieldsToAppend))
	for _, field := range fieldsToAppend {
		fmt.Printf("- %s->%s\n", field.LinkedField.LinkedModelName, field.LinkedField.Name)
	}

	// Models are setup
	fmt.Println("\nFinal Models:")
	for _, model := range models {
		fieldNames := []string{}
		for fieldName := range model.Fields {
			fieldNames = append(fieldNames, fieldName)
		}
		fmt.Printf("- %s - %v\n", model.Name, fieldNames)
	}

	// Get graphql types
	graphqlTypes := getGraphqlTypes(types)

	// Create input model types
	inputModels := inputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToInputObjects(inputModels, models)

	// Create output model types
	outputModels := outputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToOutputObjects(outputModels, models)

	// Generate query resolvers

	// Generate mutation resolvers

	// Generate final server config
	finalServerconfig := FinalServerConfig{
		EnableIntrospection: serverConfig.EnableIntrospection,
		Port:                serverConfig.Port,
		Endpoint:            serverConfig.Endpoint,
		Models:              models,
		DataLoaders:         dataLoaders,
		Types:               types,
		// Schema:
	}

	fmt.Println()
	return finalServerconfig, nil
}

func pingDataLoaders(serverConfig ServerConfig) {
}
