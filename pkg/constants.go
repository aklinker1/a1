package pkg

type contextKey string

// ContextKeyHeader is used to store the "authorization" header from the request into the graphql context
const ContextKeyHeader contextKey = "authKey"

// ChildModelDeliminator -
const ChildModelDeliminator string = "_"

// Builtin type contants
const (
	// Bool -
	Bool         string = "Bool"
	NullableBool string = "Bool?"

	// Int -
	Int         string = "Int"
	NullableInt string = "Int?"

	// Int64 -
	Int64         string = "Int64"
	NullableInt64 string = "Int64?"

	// Float -
	Float         string = "Float"
	NullableFloat string = "Float?"

	// Float64 -
	Float64         string = "Float64"
	NullableFloat64 string = "Float64?"

	// String -
	String         string = "String"
	NullableString string = "String?"

	// Date -
	Date         string = "Date"
	NullableDate string = "Date?"
)
