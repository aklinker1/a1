package new

import (
	"fmt"
	"os"

	"github.com/aklinker1/a1/pkg/utils"
	graphql "github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

// Start -
func Start(serverConfig ServerConfig) {
	// Load ENV Variables
	fmt.Println("\x1b[1mStarting Server:\x1b[0m")
	utils.Log("")

	// Load env variables
	fmt.Print("  - Loading \x1b[1m\x1b[96mEnvironment Vairables\x1b[0m from \x1b[3m.env\x1b[0m")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	isVerbose := utils.IsVerbose()
	if isVerbose {
		utils.Log("")
		utils.LogWhite("[Environment]")
		utils.Log("  - ENV_FILE: %s", envFile)
		utils.Log("  - DEV: %t", os.Getenv("DEV") == "true")
		utils.Log("  - VERBOSE: %t", isVerbose)
		utils.Log("  - STARTUP_ONLY: %t", os.Getenv("STARTUP_ONLY") == "true")
		fmt.Printf("    \x1b[92mLoaded\x1b[92m")
	}
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m")
		utils.Log("")
		fmt.Printf("        Error loading '%s': %v\n\n", envFile, err)
	} else {
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	}
	utils.Log("")

	// Parse Server Config
	fmt.Print("  - Creating the \x1b[1m\x1b[95mGraphQL Schema\x1b[0m from your models")
	finalServerconfig, errors := parseServerConfig(serverConfig)
	if len(errors) > 0 {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m")
		utils.Log("")
		fmt.Printf("Failed to create GraphQL schema, errors:\n")
		for _, err := range errors {
			fmt.Println(err)
		}
		utils.Log("")
		utils.Log("")
		os.Exit(1)
	}
	if isVerbose {
		fmt.Printf("    \x1b[92mCreated\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	utils.Log("")

	// Validate FinalServerConfig
	fmt.Print("  - Validating the \x1b[1m\x1b[93mServerConfig\x1b[0m")
	errors = validateServerConfig(finalServerconfig)
	if len(errors) > 0 {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m")
		fmt.Println()
		utils.LogRed("%d validation errors", len(errors))
		for _, err := range errors {
			fmt.Println(err)
		}
		fmt.Println()
		fmt.Println()
		os.Exit(1)
	}
	if isVerbose {
		fmt.Printf("\n    \x1b[92mValidated\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	utils.Log("")

	// Connect the data loaders
	for _, dataLoader := range finalServerconfig.DataLoaders {
		fmt.Printf("  - Connecting to \x1b[1m\x1b[94m%s\x1b[0m", dataLoader.Name)
		err := dataLoader.Connect()
		if err != nil {
			fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m")
			utils.Log("")
			utils.LogRed("Failed to create connect to the database, error: %v", err)
			utils.Log("")
			utils.Log("")
			os.Exit(1)
		}
		if isVerbose {
			fmt.Printf("\n    \x1b[92mConnected\x1b[92m")
		}
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
		utils.Log("")
	}

	// Start the server
	if os.Getenv("STARTUP_ONLY") == "true" {
		fmt.Println()
		return
	}
	startWebServer(finalServerconfig)
}

func parseServerConfig(serverConfig ServerConfig) (*FinalServerConfig, []error) {
	utils.Log("")

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
	utils.LogWhite("[Base Models - %d]", len(serverConfig.Models))
	for modelName := range serverConfig.Models {
		utils.Log("  - %s", modelName)
	}

	linkedModels := convertModelsToLinkedModels(types, serverConfig.Models)
	utils.LogWhite("[Generated Models - %d]", len(linkedModels))
	for linkedModelName, linkedModel := range linkedModels {
		utils.Log("  - %s", linkedModelName)
		serverConfig.Models[linkedModelName] = linkedModel
	}

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

	// Models are setup
	utils.LogWhite("[Final Models - %d]", len(models))
	for _, model := range models {
		fieldNames := []string{}
		for fieldName := range model.Fields {
			fieldNames = append(fieldNames, fieldName)
		}
		utils.Log("  - %s %v", model.Name, fieldNames)
	}

	// Get graphql types
	graphqlTypes := getGraphqlTypes(types)
	graphqlTypesArray := []graphql.Type{}
	for _, graphqlType := range graphqlTypes {
		graphqlTypesArray = append(graphqlTypesArray, graphqlType)
	}

	// Create input model types
	inputModels := inputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToInputObjects(inputModels, models)

	// Create output model types
	outputModels := outputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToOutputObjects(outputModels, models)

	// Generate query resolvers

	// Generate mutation resolvers

	// Creat the schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"Test": nil,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"Test": nil,
			},
		}),
		Types: graphqlTypesArray,
	})
	if err != nil {
		return nil, []error{err}
	}

	// Generate final server config
	finalServerconfig := &FinalServerConfig{
		EnableIntrospection: serverConfig.EnableIntrospection,
		Port:                serverConfig.Port,
		Endpoint:            serverConfig.Endpoint,
		Models:              models,
		DataLoaders:         dataLoaders,
		Types:               types,
		Schema:              schema,
	}

	return finalServerconfig, nil
}

func pingDataLoaders(serverConfig ServerConfig) {
}
