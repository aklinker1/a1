package pkg

import (
	graphql "github.com/graphql-go/graphql"
)

func convertModelToOutput(model Model, scalars CustomScalarMap) *graphql.Object {
	outputFields := graphql.Fields{}
	for fieldName, field := range model.Fields {
		outputFields[fieldName] = &graphql.Field{
			Type: scalars[field.Type],
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        model.Name,
		Description: model.Description,
		Fields:      outputFields,
	})
}

func convertScalarToType(scalar Scalar) graphql.Type {
	return graphql.NewScalar(graphql.ScalarConfig{
		Name:         scalar.Name,
		Description:  scalar.Description,
		Serialize:    scalar.serialize,
		ParseValue:   scalar.parse,
		ParseLiteral: scalar.parseAST,
	})
}

// Int -
var Int = Scalar{
	Name:        "Int",
	Description: "32 bit integer",
	serialize: func(value interface{}) interface{} {
		// TODO: Do this
		return value
	},
	parse: func(value interface{}) interface{} {
		// TODO: Do this
		return value
	},
	parseAST: func(valueAST ASTValue) interface{} {
		// TODO: Do this
		return valueAST
	},
}

// String -
var String = Scalar{
	Name:        "String",
	Description: "A group of characters",
	serialize: func(value interface{}) interface{} {
		// TODO: Do this
		return value
	},
	parse: func(value interface{}) interface{} {
		// TODO: Do this
		return value
	},
	parseAST: func(valueAST ASTValue) interface{} {
		// TODO: Do this
		return valueAST
	},
}
