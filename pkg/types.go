package pkg

// Int -
var Int = Scalar{
	Name:        "Int",
	Description: "32 bit integer",
	FromJSON: func(jsonValue interface{}) interface{} {
		return nil
	},
	ToJSON: func(value interface{}) interface{} {
		return nil
	},
}

// String -
var String = Scalar{
	Name:        "String",
	Description: "A group of characters",
}
