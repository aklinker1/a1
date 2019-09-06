package pkg

import (
	graphql "github.com/graphql-go/graphql"
)

// CreateSchema - Create the schema based on the server config
func CreateSchema(serverConfig ServerConfig) (graphql.Schema, error) {
	// Setup the actions
	queries := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"query": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		},
	}
	mutations := graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: graphql.Fields{
			"mutate": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		},
	}

	// Setup the GraphQL Schema
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(queries),
		Mutation: graphql.NewObject(mutations),
	}
	return graphql.NewSchema(schemaConfig)
}
