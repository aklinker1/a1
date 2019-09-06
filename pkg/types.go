package pkg

// Int -
var Int = Scalar{
	Name:        "Int",
	Description: "32 bit integer",
	FromJSON:    func(jsonValue interface{}) interface{} {},
	ToJSON:      func(value interface{}) interface{} {},
}

// String -
var String = Scalar{
	Name:     "String",
	FromJSON: func(jsonValue interface{}) interface{} {},
	ToJSON:   func(value interface{}) interface{} {},
}
