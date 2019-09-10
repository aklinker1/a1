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
	// customTypes := createScalarMap()

	// Parse Models into queries and mutations
	if isVerbose {
		fmt.Printf("    [%d Models]\n", len(serverConfig.Models))
	}
	// modelMap := createModelMap()

	// Setup the Schema
	queries := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"test": nil,
		},
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
		Name: "RootMutation",
		Fields: graphql.Fields{
			"test": nil,
		},
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
	}

	return graphql.NewSchema(schemaConfig)
}

func createScalarMap() CustomScalarMap {
	// isVerbose := os.Getenv("VERBOSE") == "true"

	return map[string]Scalar{}
}

func createModelMap() ModelMap {
	// isVerbose := os.Getenv("VERBOSE") == "true"

	return map[string]ModelMapItem{}
}
