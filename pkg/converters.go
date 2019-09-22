package pkg

import (
	"fmt"
	"os"
	"strings"

	graphql "github.com/graphql-go/graphql"
)

func (serverConfig ServerConfig) graphqlSchema() (graphql.Schema, error) {
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Println()
	}

	// Parse Custom Types
	customTypes := createTypes(serverConfig)

	// Parse Models into queries and mutations
	modelMap := createModelMap(serverConfig, customTypes)

	// Parse Queries
	queryResolvables := createQueryResolvables(modelMap)
	queryFields := graphql.Fields{}
	for _, query := range queryResolvables {
		queryFields[query.Name] = query.graphqlEntry(customTypes)
	}

	// Parse Mutations
	mutationResolvables := createMutationMap(modelMap)
	mutationFields := graphql.Fields{
		"Test": nil,
	}
	for _, mutation := range mutationResolvables {
		mutationFields[mutation.Name] = &graphql.Field{}
	}

	// Setup the Schema
	queries := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	}
	mutations := graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: mutationFields,
	}

	// Setup the GraphQL Schema
	types := []graphql.Type{}
	for _, scalar := range customTypes.Types {
		types = append(types, scalar)
	}
	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(queries),
		Mutation: graphql.NewObject(mutations),
		Types:    types,
	}

	return graphql.NewSchema(schemaConfig)
}

// Resolvers ///////////////////////////////////////////////////////////////////

func (resolver *Resolvable) graphqlResolver() func(params graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		Log("\n  %s(args: %v)", resolver.Name, resolver.Arguments)
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
		result, err := resolver.Resolver(args, fields)
		if err != nil {

			return nil, err
		}
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	}
}

func graphqlArguments(arguments []Argument, customTypes CustomTypes) graphql.FieldConfigArgument {
	config := graphql.FieldConfigArgument{}
	for _, argument := range arguments {
		config[argument.Name] = &graphql.ArgumentConfig{
			Type:         customTypes.Types[argument.Type],
			DefaultValue: argument.DefaultValue,
			Description:  argument.Description,
		}
	}
	return config
}

func (resolver *Resolvable) graphqlEntry(customTypes CustomTypes) *graphql.Field {
	return &graphql.Field{
		Name:        resolver.Name,
		Description: resolver.Description,
		Args:        graphqlArguments(resolver.Arguments, customTypes),
		Type:        customTypes.Outputs[resolver.ModelName],
		Resolve:     resolver.graphqlResolver(),
	}
}

func createQueryResolvables(modelMap ModelMap) []*Resolvable {
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

func createMutationMap(modelMap ModelMap) []Resolvable {
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

// Scalars /////////////////////////////////////////////////////////////////////

func createTypes(serverConfig ServerConfig) CustomTypes { // combine input types, output types, enums, and field types
	scalars := serverConfig.Scalars
	models := serverConfig.Models
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(len(scalars), "Scalar", "Scalars"))
	}
	scalarMap := map[string]graphql.Type{
		"String":   graphql.String,
		"Int":      graphql.Int,
		"Float":    graphql.Float,
		"Boolean":  graphql.Boolean,
		"ID":       graphql.ID,
		"DateTime": graphql.DateTime,
	}
	for _, scalar := range scalars {
		scalarMap[scalar.Name] = scalar.convertScalarToType()
	}

	outputs := map[string]*graphql.Object{}
	// Add all the basic entries for models
	for modelName, model := range models {
		outputs[modelName] = model.outputType(modelName, scalarMap)
	}
	// Basic model scalars to fill in their linked models
	for modelName, model := range models {
		for _, field := range model.Fields {
			if field.Linking != nil {
				field.Linking.graphqlType(outputs, serverConfig.Models, scalarMap, model, modelName)
			}
		}
	}

	return CustomTypes{
		Types:   scalarMap,
		Outputs: outputs,
		Inputs:  map[string]*graphql.InputObject{},
	}
}

func (scalar Scalar) convertScalarToType() graphql.Type {
	return graphql.NewScalar(graphql.ScalarConfig{
		Name:         scalar.Name,
		Description:  scalar.Description,
		Serialize:    scalar.serialize,
		ParseValue:   scalar.parse,
		ParseLiteral: scalar.parseAST,
	})
}

func (link *LinkedField) graphqlType(outputs map[string]*graphql.Object, models map[string]Model, scalars CustomScalarMap, parentModel Model, parentModelName string) {
	// Get a unique name
	var objectName string
	switch link.Type {
	case OneToOne:
		objectName = fmt.Sprintf("%s_%s", parentModelName, link.ModelName)
	case OneToMany:
		var pluralAddition = "s"
		if strings.HasSuffix(link.ModelName, "s") {
			pluralAddition = ""
		}
		objectName = fmt.Sprintf("%s_%s%s", parentModelName, link.ModelName, pluralAddition)
	}

	// Get the fields
	linkModel := models[link.ModelName]
	outputFields := graphql.Fields{}
	for fieldName, field := range linkModel.Fields {
		if field.Linking != nil {
			field.Linking.graphqlType(outputs, models, scalars, linkModel, link.ModelName)
		}
		outputFields[fieldName] = &graphql.Field{
			Name: fieldName,
			Type: scalars[field.Type],
		}
	}
	// Add the new linked object to the outputs
	object := graphql.NewObject(graphql.ObjectConfig{
		Name:        objectName,
		Description: linkModel.Description,
		Fields:      outputFields,
	})
	outputs[objectName] = object

	// Add the field to the existing output object
	outputs[parentModelName].AddFieldConfig(link.AccessedAs, &graphql.Field{
		Name: link.AccessedAs,
		Type: object,
	})

	// Add the reverse field to the other existing output object
	var reverseAccessType graphql.Type
	if link.Type == OneToMany {
		reverseAccessType = graphql.NewList(outputs[parentModelName])
	} else {
		reverseAccessType = outputs[parentModelName]
	}
	outputs[link.ModelName].AddFieldConfig(link.ReverseAccessedAs, &graphql.Field{
		Name: link.ReverseAccessedAs,
		Type: reverseAccessType,
	})
}

// Models //////////////////////////////////////////////////////////////////////

func (model Model) outputType(modelName string, scalars CustomScalarMap) *graphql.Object {
	outputFields := graphql.Fields{}
	for fieldName, field := range model.Fields {
		outputFields[fieldName] = &graphql.Field{
			Name: fieldName,
			Type: scalars[field.Type],
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        modelName,
		Description: model.Description,
		Fields:      outputFields,
	})
}

func (model Model) createModelItem(modelName string, serverconfig ServerConfig) ModelMapItem {
	queries := []*Resolvable{
		selectOneQuery(modelName, model, serverconfig),
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

func createModelMap(serverConfig ServerConfig, customTypes CustomTypes) ModelMap {
	models := serverConfig.Models
	isVerbose := os.Getenv("VERBOSE") == "true"
	if isVerbose {
		fmt.Printf("    [%s]\n", pluralize(len(models), "Model", "Models"))
	}

	modelMap := map[string]ModelMapItem{}

	for modelName, model := range models {
		if isVerbose {
			fmt.Printf("      - %s\n", modelName)
		}
		modelMap[modelName] = model.createModelItem(modelName, serverConfig)
	}

	return modelMap
}
