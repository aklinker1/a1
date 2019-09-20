package pkg

import (
	"fmt"
	"os"

	graphql "github.com/graphql-go/graphql"
)

// CreateSchema - Create the schema based on the server config
func CreateSchema(serverConfig ServerConfig) (graphql.Schema, error) {
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Println()
	}

	// Parse Custom Types
	scalarMap := createScalarMap(serverConfig)
	scalars := []graphql.Type{}
	for _, scalar := range scalarMap {
		scalars = append(scalars, scalar)
	}

	// Parse Models into queries and mutations
	modelMap := createModelMap(serverConfig, scalarMap)

	// Parse Queries
	queryResolvables := createQueryResolvables(modelMap, scalarMap)
	queryFields := graphql.Fields{}
	for _, query := range queryResolvables {
		queryFields[query.Name] = convertResolver(query, scalarMap)
	}

	// Parse Mutations
	// mutationResolvables := createMutationMap(modelMap)
	mutationFields := graphql.Fields{
		"Test": nil,
	}
	// for _, mutation := range mutationResolvables {
	// 	mutationFields[mutation.Name] = &graphql.Field{}
	// }

	// Setup the Schema
	queries := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
		// Fields: graphql.Fields{
		// 	"query": &graphql.Field{
		// 		Type: graphql.String,
		// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// 			return "world", nil
		// 		},
		// 	},
		// },
	}
	mutations := graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: mutationFields,
		// Fields: graphql.Fields{
		// 	"mutate": &graphql.Field{
		// 		Type: graphql.String,
		// 		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// 			return "world", nil
		// 		},
		// 	},
		// },
	}

	// Setup the GraphQL Schema
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(queries),
		Mutation: graphql.NewObject(mutations),
		Types:    scalars,
	}

	return graphql.NewSchema(schemaConfig)
}

// Resolvers ///////////////////////////////////////////////////////////////////

func middlewareResolver(resolver *Resolvable) func(params graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
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
		fields := StringMap{}

		// Call Resolver
		return resolver.Resolver(args, fields)
	}
}

func convertArguments(arguments []Argument, scalars map[string]graphql.Type) graphql.FieldConfigArgument {
	config := graphql.FieldConfigArgument{}
	for _, argument := range arguments {
		config[argument.Name] = &graphql.ArgumentConfig{
			Type:         scalars[argument.Type],
			DefaultValue: argument.DefaultValue,
			Description:  argument.Description,
		}
	}
	return config
}

func convertResolver(resolver *Resolvable, scalars map[string]graphql.Type) *graphql.Field {
	return &graphql.Field{
		Name:        resolver.Name,
		Description: resolver.Description,
		Args:        convertArguments(resolver.Arguments, scalars),
		Type:        convertModelToOutput(resolver.Model, scalars),
		Resolve:     middlewareResolver(resolver),
	}
}

// Scalars /////////////////////////////////////////////////////////////////////

func createScalarMap(serverConfig ServerConfig) CustomScalarMap {
	scalars := serverConfig.Scalars
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(len(scalars), "Scalar", "Scalars"))
	}
	scalarMap := map[string]graphql.Type{
		"String": graphql.String,
		"Int":    graphql.Int,
	}
	for _, scalar := range scalars {
		t := convertScalarToType(scalar)
		scalarMap[scalar.Name] = t
	}

	return scalarMap
}

// Modes ///////////////////////////////////////////////////////////////////////

func createModelItem(model Model, databaseDriver DatabaseDriver) ModelMapItem {
	queries := []*Resolvable{
		selectOneQuery(model, databaseDriver),
	}
	if model.GraphQL.CustomQueries != nil {
		queries = append(queries, model.GraphQL.CustomQueries...)
	}
	mutations := []*Resolvable{}
	return ModelMapItem{
		Model:     model,
		Queries:   append(queries, model.GraphQL.CustomQueries...),
		Mutations: append(mutations, model.GraphQL.CustomMutations...),
	}
}

func createModelMap(serverConfig ServerConfig, scalars CustomScalarMap) ModelMap {
	models := serverConfig.Models
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(len(models), "Model", "Models"))
	}

	modelMap := map[string]ModelMapItem{}

	for _, model := range models {
		if isVerbose {
			fmt.Printf("      - %s\n", model.Name)
		}
		modelMap[model.Name] = createModelItem(model, serverConfig.DatabaseDriver)
	}

	return modelMap
}

// Queries /////////////////////////////////////////////////////////////////////

func createQueryResolvables(modelMap ModelMap, scalars CustomScalarMap) []*Resolvable {
	isVerbose := os.Getenv("VERBOSE") == "true"
	queryCount := 0
	for _, modelItem := range modelMap {
		queryCount += len(modelItem.Queries)
	}
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(queryCount, "Query", "Queries"))
	}
	queries := []*Resolvable{}
	for _, modelItem := range modelMap {
		queries = append(queries, modelItem.Queries...)
	}
	if isVerbose {
		for _, query := range queries {
			fmt.Printf("      - %s\n", query.Name)
		}
	}
	return queries
}

// Mutations ///////////////////////////////////////////////////////////////////

func createMutationMap(modelMap ModelMap, scalars CustomScalarMap) []Resolvable {
	isVerbose := os.Getenv("VERBOSE") == "true"
	mutationCount := 0
	for _, value := range modelMap {
		mutationCount += len(value.Mutations)
	}
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(mutationCount, "Mutation", "Mutations"))
	}
	return nil
}
