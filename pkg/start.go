package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/aklinker1/a1/pkg/utils"
	graphql "github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

// Start -
func Start(serverConfig ServerConfig) {
	fmt.Println("\x1b[1mStarting Server:\x1b[0m")
	utils.Log("")

	// Load ENV Variables
	envFile := os.Getenv("ENV_FILE")
	fmt.Printf("  - Loading \x1b[1m\x1b[96mEnvironment Vairables\x1b[0m from \x1b[3m%s\x1b[0m", envFile)
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	isVerbose := utils.IsVerbose()
	if isVerbose {
		utils.Log("")
		utils.LogWhite("[a1 Environment]")
		utils.Log("  - ENV_FILE: %s", envFile)
		utils.Log("  - DEV: %t", os.Getenv("DEV") == "true")
		utils.Log("  - VERBOSE: %t", isVerbose)
		utils.Log("  - STARTUP_ONLY: %t", os.Getenv("STARTUP_ONLY") == "true")
		utils.LogWhite("[Custom Environment]")
		printOtherENV(envFile)
		fmt.Printf("    \x1b[92mLoaded\x1b[92m")
	}
	if err != nil {
		fmt.Println(" \x1b[91m\x1b[1m(✘)\x1b[0m")
		utils.Log("")
		fmt.Printf("      Error loading '%s':\n", envFile)
		fmt.Printf("      \x1b[91m%v\x1b[0m\n\n", err)
	} else {
		fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	}
	utils.Log("")

	// Parse Server Config
	fmt.Print("  - Parsing \x1b[1m\x1b[93mServerConfig\x1b[0m input")
	finalServerConfig, errors := parseServerConfig(serverConfig)
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
	errors = validateServerConfig(finalServerConfig)
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

	// Create the GraphQL Schema
	fmt.Print("  - Creating the \x1b[1m\x1b[95mGraphQL Schema\x1b[0m from your models")
	schema, err := createSchema(finalServerConfig)
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
	finalServerConfig.GraphQLSchema = schema
	if isVerbose {
		fmt.Printf("\n    \x1b[92mCreated\x1b[92m")
	}
	fmt.Println(" \x1b[92m\x1b[1m(✔)\x1b[0m")
	utils.Log("")

	// Connect the data loaders
	for _, dataLoader := range finalServerConfig.DataLoaders {
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
	startWebServer(finalServerConfig)
}

func printOtherENV(envFile string) {
	bytes, _ := ioutil.ReadFile(envFile)
	fileContent := string(bytes)
	lines := strings.Split(fileContent, "\n")
	ignoredKeys := []string{"DEV", "VERBOSE", "STARTUP_ONLY"}
	re := regexp.MustCompile(`(?m)(.*)=(.*)`)

	for _, line := range lines {
		isIgnored := false
		for _, ignoredKey := range ignoredKeys {
			if strings.HasPrefix(line, ignoredKey) {
				isIgnored = true
			}
		}
		if !isIgnored {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 0 {
				key := matches[1]
				value := matches[2]
				utils.Log("  - %s: %s", key, value)
			}
		}
	}
}

func parseServerConfig(serverConfig ServerConfig) (*FinalServerConfig, []error) {
	utils.Log("")
	baseModelMap := ModelMap{}
	for modelName, model := range serverConfig.Models {
		baseModelMap[modelName] = model
	}

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
		model := convertModelToFinalWithoutLinkedFields(dataLoaders, types, modelName, model)
		models[modelName] = model
	}

	// Copy model properties to extended models
	extendModels(models)
	baseModels := FinalModelMap{}
	for modelName, model := range models {
		if _, ok := baseModelMap[modelName]; ok {
			baseModels[modelName] = model
		}
	}

	// Append linked fields to models
	fieldsToAppend := getLinkedFieldsToAppend(serverConfig.Models, models)
	for _, modelAndField := range fieldsToAppend {
		modelAndField.Model.Fields[modelAndField.LinkedField.Name] = modelAndField.LinkedField
	}

	// Models are setup
	if utils.IsVerbose() {
		utils.LogWhite("[Final Models - %d]", len(models))
		for _, model := range models {
			fieldNames := []string{}
			for fieldName := range model.Fields {
				fieldNames = append(fieldNames, fieldName)
			}
			sort.Strings(fieldNames)
			utils.Log("  - %s %v", model.Name, fieldNames)
		}
	}

	// Get graphql types
	graphqlTypes := getGraphqlTypes(types)
	graphqlTypesArray := []graphql.Type{}
	for _, graphqlType := range graphqlTypes {
		graphqlTypesArray = append(graphqlTypesArray, graphqlType)
	}

	// Create output model types
	outputModels := outputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToOutputObjects(outputModels, models)

	// Create input model types
	inputModels := inputModelsWithoutLinkedFields(graphqlTypes, models)
	addLinksToInputObjects(inputModels, models)

	// Combine Types
	allTypes := graphql.TypeMap{}
	for graphqlTypeName, graphqlType := range graphqlTypes {
		allTypes[graphqlTypeName] = graphqlType
	}
	for inputTypeName, inputType := range inputModels {
		allTypes[inputTypeName] = inputType
	}
	for outputTypeName, outputType := range outputModels {
		allTypes[outputTypeName] = outputType
	}

	// Generate final server config
	finalServerConfig := &FinalServerConfig{
		EnableIntrospection: serverConfig.EnableIntrospection,
		Port:                serverConfig.Port,
		Endpoint:            serverConfig.Endpoint,
		Models:              models,
		DataLoaders:         dataLoaders,
		Types:               types,
	}

	// Generate query resolvers
	queryResolvables := []Resolvable{}
	for _, model := range baseModels {
		queryResolvables = append(queryResolvables, generateQueriesForModel(finalServerConfig, model)...)
	}
	queries := graphql.Fields{}
	utils.LogWhite("[Queries - %d]", len(queryResolvables))
	for _, query := range queryResolvables {
		queries[query.Name] = query.graphqlResolverEntry(allTypes)
		utils.Log("  - %s", query.Name)
	}
	finalServerConfig.GraphQLQueries = queries

	// Generate mutation resolvers
	finalServerConfig.GraphQLMutations = graphql.Fields{
		"test": nil,
	}

	return finalServerConfig, nil
}

func createSchema(serverConfig *FinalServerConfig) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: serverConfig.GraphQLQueries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootMutation",
			Fields: serverConfig.GraphQLMutations,
		}),
		Types: serverConfig.GraphqlTypes,
	})
}
