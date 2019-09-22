package pkg

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
